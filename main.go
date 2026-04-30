package main

import (
	"go-url-shortener/database"
	"go-url-shortener/handlers"
	"go-url-shortener/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.URL{}, &models.User{})

	r := gin.Default()

	// Auth routes (публичные)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected dashboard
	r.GET("/dashboard", handlers.AuthMiddleware(), func(c *gin.Context) {
		user, _ := c.Get("username")
		c.JSON(200, gin.H{"message": "Welcome, " + user.(string)})
	})
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
	r.GET("/health/all", handlers.CheckAllURLsHealth)
	r.GET("/health/:id", handlers.CheckURLHealth)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
