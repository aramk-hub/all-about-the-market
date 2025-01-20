package config

import (
    "log"
    "os"
)

type AppConfig struct {
    ClientID      string
    CognitoDomain string
    UserPoolID    string
    JWKSURL       string
    RedirectURI   string
    FrontendURL   string
}

func LoadEnvVariables() *AppConfig {
    clientID := os.Getenv("COGNITO_CLIENT_ID")
    cognitoDomain := os.Getenv("COGNITO_DOMAIN")
    userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
    jwksURL := os.Getenv("JWKS_URL")
    redirectURI := os.Getenv("COGNITO_REDIRECT_URI")
    frontendURL := os.Getenv("FRONTEND_REDIRECT_URL")

    // Validate that all required environment variables are present
    if clientID == "" || cognitoDomain == "" || userPoolID == "" || jwksURL == "" || redirectURI == "" || frontendURL == "" {
        log.Fatalf("Required environment variables are missing! Please check your .env file or environment settings.")
    }

    log.Printf("Cognito Client ID: %s", clientID)
    log.Printf("Cognito Domain: %s", cognitoDomain)
    log.Printf("User Pool ID: %s", userPoolID)
    log.Printf("JWKS URL: %s", jwksURL)
    log.Printf("Redirect URI: %s", redirectURI)
    log.Printf("Frontend Redirect URL: %s", frontendURL)

    return &AppConfig{
        ClientID:      clientID,
        CognitoDomain: cognitoDomain,
        UserPoolID:    userPoolID,
        JWKSURL:       jwksURL,
        RedirectURI:   redirectURI,
        FrontendURL:   frontendURL,
    }
}
