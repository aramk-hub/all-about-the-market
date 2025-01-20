package handlers

import (
	"github.com/gin-gonic/gin"
	"all-about-the-market/backend/models"
	"all-about-the-market/backend/database"
	"net/http"
)

// CreateStock handles POST request to add a stock to a portfolio
func CreateStock(c *gin.Context) {
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Save stock to database
	if err := database.DB.Create(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock"})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

// GetStocks handles GET request to retrieve all stocks in a portfolio
func GetStocks(c *gin.Context) {
	var stocks []models.Stock
	if err := database.DB.Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stocks"})
		return
	}
	c.JSON(http.StatusOK, stocks)
}
