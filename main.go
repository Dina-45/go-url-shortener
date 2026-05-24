package main

import (
	"go-url-shortener/database"
	"go-url-shortener/handlers"
	"go-url-shortener/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.URL{}, &models.User{})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/short/:code", handlers.RedirectShortURL)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := r.Group("/")
	auth.Use(handlers.AuthMiddleware())
	{
		auth.GET("/dashboard", func(c *gin.Context) {
			user, _ := c.Get("username")
			c.JSON(200, gin.H{"message": "Welcome, " + user.(string)})
		})

		auth.GET("/urls", handlers.GetAllURLs)
		auth.GET("/urls/:id", handlers.GetURLByID)
		auth.POST("/urls", handlers.CreateURL)
		auth.PUT("/urls/:id", handlers.UpdateURL)
		auth.DELETE("/urls/:id", handlers.DeleteURL)
		auth.GET("/count/:id", handlers.GetClicks)
		auth.POST("/urls/bulk", handlers.CreateBulkURLs)
		auth.DELETE("/urls", handlers.DeleteAllURLs)
		auth.GET("/stats", handlers.GetStats)
		auth.GET("/health/all", handlers.CheckAllURLsHealth)
		auth.GET("/health/:id", handlers.CheckURLHealth)
	}
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	r.Run(":8080")
}
