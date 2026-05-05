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
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/middlewares" 
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

	// TAMBAHAN: Dependency injection untuk Menu
	menuRepo := repositories.NewMenuRepository(db)
	menuService := services.NewMenuService(menuRepo)
	menuHandler := handlers.NewMenuHandler(menuService)

	// TAMBAHAN: Dependency injection untuk Order
	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo, menuRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	r := gin.Default()

	// route ping for connection test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Connected successfully!",
		})
	})

	// route api tenant
	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/orders", orderHandler.Create)
		api.GET("/admin/orders", orderHandler.GetAllOrders)
        api.PUT("/admin/orders/:id/pay", orderHandler.MarkOrderAsPaid)
	}

	// TAMBAHAN: Protected routes (HANYA bisa diakses DENGAN token JWT)
	protected := r.Group("/api/v1")
	// Pastikan AuthMiddleware kamu menerima secretKey jika desain kodemu membutuhkannya
	// Ubah menjadi middleware.AuthMiddleware() jika tidak butuh parameter secretKey
	protected.Use(middlewares.AuthMiddleware(secretKey)) 
	{
		protected.POST("/menus", menuHandler.Create)       // Create menu baru
		protected.GET("/menus", menuHandler.GetAll)        // Lihat semua menu milik tenant tersebut
		protected.PUT("/menus/:id", menuHandler.Update)    // Update menu berdasarkan ID
		protected.DELETE("/menus/:id", menuHandler.Delete) // Hapus menu berdasarkan ID
		protected.GET("/tenant/orders", orderHandler.GetTenantOrders)
        protected.PUT("/tenant/orders/:id/status", orderHandler.UpdateTenantOrderStatus)
	}

	// port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}