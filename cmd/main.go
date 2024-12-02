package main

import (
	"all-about-the-market/config"
	"all-about-the-market/database"
	"all-about-the-market/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.GetUserFromCognito("aramkazorian@gmail.com")
	// Connect to PostgreSQL database
	database.ConnectDB()

	// Initialize Cognito client
	config.InitCognito()

	// Set up Gin router
	r := gin.Default()

	// // Define routes for Portfolio
	// r.GET("/portfolios", handlers.GetPortfolios)
	// r.POST("/portfolios", handlers.CreatePortfolio)

	// // Define routes for Stock
	// r.GET("/portfolios/:portfolio_id/stocks", handlers.GetStocks)
	// r.POST("/portfolios/:portfolio_id/stocks", handlers.CreateStock)

	// Start server
	r.Run(":8080")
}
