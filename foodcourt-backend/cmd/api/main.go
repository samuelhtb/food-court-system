package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/config" 
)

func main() {
	config.ConnectDB()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Connected successfully!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}