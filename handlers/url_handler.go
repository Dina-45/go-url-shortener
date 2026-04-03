package handlers

import (
	"go-url-shortener/database"
	"go-url-shortener/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /urls
func GetAllURLs(c *gin.Context) {
	var urls []models.URL
	database.DB.Find(&urls)
	c.JSON(http.StatusOK, urls)
}

// GET /urls/:id
func GetURLByID(c *gin.Context) {
	id := c.Param("id")
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, url)
}

// POST /urls
func CreateURL(c *gin.Context) {
	var url models.URL
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&url)
	c.JSON(http.StatusCreated, url)
}

// PUT /urls/:id
func UpdateURL(c *gin.Context) {
	id := c.Param("id")
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	var input models.URL
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Model(&url).Updates(models.URL{OriginalURL: input.OriginalURL, ShortCode: input.ShortCode})
	c.JSON(http.StatusOK, url)
}

// DELETE /urls/:id
func DeleteURL(c *gin.Context) {
	id := c.Param("id")
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	database.DB.Delete(&url)
	c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
}

// GET /short/:code
func RedirectShortURL(c *gin.Context) {
	code := c.Param("code")
	var url models.URL
	if err := database.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"original_url": url.OriginalURL})
}

// GET /count/:id
func GetClicks(c *gin.Context) {
	id := c.Param("id")
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"clicks": 0}) // пока заглушка
}

// POST /urls/bulk
func CreateBulkURLs(c *gin.Context) {
	var inputs []models.URL
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, input := range inputs {
		database.DB.Create(&models.URL{OriginalURL: input.OriginalURL, ShortCode: input.ShortCode})
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Bulk URLs created"})
}

// DELETE /urls
func DeleteAllURLs(c *gin.Context) {
	database.DB.Where("1 = 1").Delete(&models.URL{})
	c.JSON(http.StatusOK, gin.H{"message": "All URLs deleted"})
}

// GET /stats
func GetStats(c *gin.Context) {
	var count int64
	database.DB.Model(&models.URL{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"total_urls": count})
}
