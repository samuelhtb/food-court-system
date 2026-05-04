package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/config"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/handlers"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment default")
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("Error: JWT_SECRET belum diatur di file .env")
	}

	// db connection
	db := config.ConnectDB()

	// dependency injection
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, secretKey)

	r := gin.Default()

	// route ping for connection test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Connected successfully!",
		})
	})

	// route api
	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}