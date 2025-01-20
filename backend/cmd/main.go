package main

import (
	"log"

	"all-about-the-market/backend/config"
	"all-about-the-market/backend/database"
	"all-about-the-market/backend/handlers"
	"all-about-the-market/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load and validate application configuration
	appConfig := config.LoadEnvVariables()

	// Connect to PostgreSQL database
	database.ConnectDB()

	// Initialize Cognito client
	config.InitCognito()

	// Fetch JWKS keys for token validation
	err = utils.FetchJWKS()
	if err != nil {
		log.Fatalf("Failed to fetch JWKS: %v", err)
	}

	// Example function to fetch user (ensure Cognito is initialized first)
	handlers.GetUserFromCognito("aramkazorian@gmail.com")

	// Set up Gin router
	r := gin.Default()

	// Define route for Cognito authentication callback
	r.GET("/auth/callback", func(c *gin.Context) {
		handlers.AuthCallbackHandler(c.Writer, c.Request, appConfig)
	})

	// Start server
	log.Printf("Server started on :8080 with Cognito Domain: %s", appConfig.CognitoDomain)
	r.Run(":8080")
}
