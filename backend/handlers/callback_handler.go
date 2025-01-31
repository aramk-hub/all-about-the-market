package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"all-about-the-market/backend/config"
	"all-about-the-market/backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthCallbackHandler(c *gin.Context) {
	// Extract authorization code from query parameters
	code := c.Query("code")
	if code == "" {
		log.Println("Authorization code not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	// Exchange the authorization code for tokens
	appConfig := config.LoadEnvVariables() // Replace with your app's way of accessing AppConfig
	tokens, err := exchangeCodeForTokensWithoutSecret(code, appConfig)
	if err != nil {
		log.Printf("Failed to exchange code for tokens: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for tokens"})
		return
	}

	// Extract ID token and access token
	idToken, ok := tokens["id_token"].(string)
	if !ok {
		log.Println("ID token missing from Cognito response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID token missing"})
		return
	}

	// Validate the ID token
	if _, err := utils.ValidateToken(idToken); err != nil {
		log.Printf("Failed to validate ID token: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
		return
	}

	accessToken, ok := tokens["access_token"].(string)
	if !ok {
		log.Println("Access token missing from Cognito response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Access token missing"})
		return
	}

	// Set tokens as secure, HTTP-only cookies
	c.SetCookie("id_token", idToken, 3600, "/", "", false, true)
	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)

	// Redirect the user to the frontend dashboard
	c.Redirect(http.StatusSeeOther, "http://localhost:3000/dashboard")
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
