package handlers

import (
	"net/http"
	"time"

	"go-url-shortener/database"
	"go-url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var restyClient = resty.New().
	SetTimeout(5 * time.Second)

func checkURL(original string) bool {
	resp, err := restyClient.R().Head(original)
	if err != nil {
		return false
	}
	return resp.StatusCode() < 400
}

func CheckURLHealth(c *gin.Context) {
	var url models.URL
	if err := database.DB.First(&url, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL не найден"})
		return
	}

	alive := checkURL(url.OriginalURL)
	database.DB.Model(&url).Update("is_alive", alive)

	status := "alive"
	if !alive {
		status = "dead"
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           url.ID,
		"original_url": url.OriginalURL,
		"status":       status,
		"is_alive":     alive,
	})
}

func CheckAllURLsHealth(c *gin.Context) {
	var urls []models.URL
	database.DB.Find(&urls)

	type Result struct {
		ID          uint   `json:"id"`
		OriginalURL string `json:"original_url"`
		Status      string `json:"status"`
	}

	var results []Result

	for _, u := range urls {
		alive := checkURL(u.OriginalURL)
		database.DB.Model(&u).Update("is_alive", alive)

		status := "alive"
		if !alive {
			status = "dead"
		}
		results = append(results, Result{
			ID:          u.ID,
			OriginalURL: u.OriginalURL,
			Status:      status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"checked": len(results),
		"results": results,
	})
}
