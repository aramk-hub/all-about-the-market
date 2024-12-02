package handlers

import (
	"github.com/gin-gonic/gin"
	"all-about-the-market/models"
	"all-about-the-market/database"
	"net/http"
)

// CreatePortfolio handles POST request to create a new portfolio
func CreatePortfolio(c *gin.Context) {
	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Save portfolio to database
	if err := database.DB.Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// GetPortfolios handles GET request to retrieve all portfolios
func GetPortfolios(c *gin.Context) {
	var portfolios []models.Portfolio
	if err := database.DB.Find(&portfolios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch portfolios"})
		return
	}
	c.JSON(http.StatusOK, portfolios)
}
