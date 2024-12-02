package main

import (
	"github.com/gin-gonic/gin"
	"all-about-the-market/database"
	"all-about-the-market/handlers"
)

func main() {
	// Connect to PostgreSQL database
	database.ConnectDB()

	// Set up Gin router
	r := gin.Default()

	// Define routes for Portfolio
	r.GET("/portfolios", handlers.GetPortfolios)
	r.POST("/portfolios", handlers.CreatePortfolio)

	// Define routes for Stock
	r.GET("/portfolios/:portfolio_id/stocks", handlers.GetStocks)
	r.POST("/portfolios/:portfolio_id/stocks", handlers.CreateStock)

	// Start server
	r.Run(":8080")
}
