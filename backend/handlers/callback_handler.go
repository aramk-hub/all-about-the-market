package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"all-about-the-market/backend/utils"
	"all-about-the-market/backend/config"
)

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request, appConfig *config.AppConfig) {
	// Extract authorization code from query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Authorization code not found")
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}
	log.Printf("Authorization code received: %s\n", code)

	// Exchange the authorization code for tokens
	tokens, err := exchangeCodeForTokensWithoutSecret(code, appConfig)
	if err != nil {
		log.Printf("Failed to exchange code for tokens: %v\n", err)
		http.Error(w, "Failed to exchange code for tokens", http.StatusInternalServerError)
		return
	}

	// Extract ID token
	idToken, ok := tokens["id_token"].(string)
	if !ok {
		log.Println("ID token missing from Cognito response")
		http.Error(w, "ID token missing", http.StatusInternalServerError)
		return
	}

	// Validate the ID token
	if _, err := utils.ValidateToken(idToken); err != nil {
		log.Printf("Failed to validate ID token: %v\n", err)
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	// Log the tokens for debugging (remove in production)
	log.Printf("Tokens received: %+v\n", tokens)

	// Redirect to the frontend application after successful authentication
	frontendURL := os.Getenv("FRONTEND_REDIRECT_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000/dashboard" // Default fallback
	}
	http.Redirect(w, r, frontendURL, http.StatusSeeOther)
}

// exchangeCodeForTokensWithoutSecret exchanges the authorization code for tokens without using a client secret
func exchangeCodeForTokensWithoutSecret(code string, appConfig *config.AppConfig) (map[string]interface{}, error) {
	tokenEndpoint := appConfig.CognitoDomain + "/oauth2/token"

	// Prepare request payload
	data := strings.NewReader("grant_type=authorization_code" +
		"&client_id=" + appConfig.ClientID +
		"&redirect_uri=" + appConfig.RedirectURI +
		"&code=" + code)

	req, err := http.NewRequest("POST", tokenEndpoint, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from Cognito: %s", resp.Status)
		return nil, err
	}

	var tokens map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}
