package main

import (
	"go-url-shortener/database"
	"go-url-shortener/handlers"
	"go-url-shortener/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.URL{})

	r := gin.Default()

	r.GET("/urls", handlers.GetAllURLs)
	r.GET("/urls/:id", handlers.GetURLByID)
	r.POST("/urls", handlers.CreateURL)
	r.PUT("/urls/:id", handlers.UpdateURL)
	r.DELETE("/urls/:id", handlers.DeleteURL)
	r.GET("/short/:code", handlers.RedirectShortURL)
	r.GET("/count/:id", handlers.GetClicks)
	r.POST("/urls/bulk", handlers.CreateBulkURLs)
	r.DELETE("/urls", handlers.DeleteAllURLs)
	r.GET("/stats", handlers.GetStats)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
