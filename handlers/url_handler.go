package handlers

import (
	"go-url-shortener/database"
	"go-url-shortener/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func generateUniqueCode() string {
	for {
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		code := string(b)

		var existing models.URL
		if err := database.DB.Where("short_code = ?", code).First(&existing).Error; err != nil {
			// Не нашли — код свободен
			return code
		}
	}
}

func GetAllURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var urls []models.URL
	database.DB.Where("user_id = ?", userID).Find(&urls)
	c.JSON(http.StatusOK, urls)
}

func GetURLByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	if url.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, url)
}

func CreateURL(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		OriginalURL string `json:"original_url" binding:"required"`
		ShortCode   string `json:"short_code"`
		ExpiresIn   *int   `json:"expires_in"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortCode := input.ShortCode
	if shortCode == "" || shortCode == " " {
		shortCode = generateUniqueCode()
	} else {
		var existing models.URL
		if err := database.DB.Where("short_code = ?", shortCode).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Short code already taken"})
			return
		}
	}

	url := models.URL{
		OriginalURL: input.OriginalURL,
		ShortCode:   shortCode,
		UserID:      userID.(uint),
	}

	if input.ExpiresIn != nil && *input.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(*input.ExpiresIn) * time.Hour)
		url.ExpiresAt = &t
	}

	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать"})
		return
	}

	host := c.Request.Host
	c.JSON(http.StatusCreated, gin.H{
		"id":           url.ID,
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"short_url":    "http://" + host + "/short/" + url.ShortCode,
		"clicks":       url.Clicks,
		"expires_at":   url.ExpiresAt,
		"created_at":   url.CreatedAt,
	})
}

func UpdateURL(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	if url.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
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

func DeleteURL(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	if url.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	database.DB.Delete(&url)
	c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
}

func RedirectShortURL(c *gin.Context) {
	code := c.Param("code")
	var url models.URL
	if err := database.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
		return
	}

	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Link has expired"})
		return
	}

	database.DB.Model(&url).UpdateColumn("clicks", url.Clicks+1)

	c.Redirect(http.StatusFound, url.OriginalURL)
}

func GetClicks(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	if url.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     url.ID,
		"clicks": url.Clicks,
	})
}

func CreateBulkURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var inputs []struct {
		OriginalURL string `json:"original_url"`
		ShortCode   string `json:"short_code"`
	}
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var created []models.URL
	for _, input := range inputs {
		code := input.ShortCode
		if code == "" {
			code = generateUniqueCode()
		}
		u := models.URL{
			OriginalURL: input.OriginalURL,
			ShortCode:   code,
			UserID:      userID.(uint),
		}
		if err := database.DB.Create(&u).Error; err == nil {
			created = append(created, u)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bulk URLs created",
		"created": len(created),
		"urls":    created,
	})
}

func DeleteAllURLs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	database.DB.Where("user_id = ?", userID).Delete(&models.URL{})
	c.JSON(http.StatusOK, gin.H{"message": "All your URLs deleted"})
}

func GetStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var totalCount int64
	database.DB.Model(&models.URL{}).Where("user_id = ?", userID).Count(&totalCount)

	var totalClicks int64
	database.DB.Model(&models.URL{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(clicks), 0)").
		Scan(&totalClicks)

	var expiredCount int64
	database.DB.Model(&models.URL{}).
		Where("user_id = ? AND expires_at IS NOT NULL AND expires_at < ?", userID, time.Now()).
		Count(&expiredCount)

	c.JSON(http.StatusOK, gin.H{
		"total_urls":   totalCount,
		"total_clicks": totalClicks,
		"expired_urls": expiredCount,
	})
}
