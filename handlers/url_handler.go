package handlers

import (
	"net/http"
	"urlshortener/database"
	"urlshortener/models"

	"github.com/gin-gonic/gin"
)

func CreateURL(c *gin.Context) {
	var url models.URL
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&url)
	c.JSON(http.StatusOK, url)
}

func GetAllURLs(c *gin.Context) {
	var urls []models.URL
	database.DB.Find(&urls)
	c.JSON(http.StatusOK, gin.H{"data": urls})
}
