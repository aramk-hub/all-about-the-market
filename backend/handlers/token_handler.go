package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenHandler(c *gin.Context) {
	// Retrieve cookies
	idToken, err := c.Cookie("id_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID token not found"})
		return
	}

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
		return
	}

	// Return tokens as JSON
	c.JSON(http.StatusOK, gin.H{
		"id_token":     idToken,
		"access_token": accessToken,
	})
}
