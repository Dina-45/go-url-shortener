package main

import (
	"go-url-shortener/database"
	"go-url-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()


	r.POST("/urls", handlers.CreateURL)       // 1
	r.GET("/urls", handlers.GetAllURLs)       // 2
	r.GET("/urls/:id", handlers.GetURL)       // 3
	r.PUT("/urls/:id", handlers.UpdateURL)    // 4
	r.DELETE("/urls/:id", handlers.DeleteURL) // 5
	r.GET("/r/:code", handlers.RedirectURL)   // 6
	r.GET("/stats/:code", handlers.GetStats)  // 7

	r.GET("/urls/search", handlers.SearchURLs)
	r.GET("/urls/top", handlers.TopURLs)
	r.DELETE("/urls", handlers.DeleteAllURLs)

	r.Run(":8080")
}