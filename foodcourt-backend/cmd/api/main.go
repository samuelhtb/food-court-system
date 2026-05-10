package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/config"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/handlers"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/middlewares"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
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

	// migration
	err = db.AutoMigrate(&models.User{}, &models.Menu{}, &models.Order{}, &models.SubOrder{}, &models.OrderItem{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi: %v", err)
	}

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

	// TAMBAHAN: Midtrans
	midtransCfg := config.LoadMidtransConfig()
	midtransService := services.NewMidtransService(orderRepo, midtransCfg)
	midtransHandler := handlers.NewMidtransHandler(midtransService)

	orderHandler := handlers.NewOrderHandler(orderService, midtransService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Mengizinkan Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

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
		api.POST("/register-admin", userHandler.RegisterAdmin)
		api.POST("/login", userHandler.Login)

		// Pelanggan bisa melihat menu dan membuat pesanan tanpa perlu login (opsional, sesuaikan bisnis)
		api.GET("/menus", menuHandler.GetAllPublic)
		api.POST("/orders", orderHandler.Create)

		api.GET("/orders/:id", orderHandler.GetOrderDetails)

		// Midtrans Webhooks & Verification
		api.POST("/midtrans/notification", midtransHandler.HandleNotification)
		api.POST("/midtrans/verify", midtransHandler.VerifyLocalPayment)
	}

	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware(secretKey))
	{
		// Tenant
		tenantRoutes := protected.Group("/")
		tenantRoutes.Use(middlewares.RoleMiddleware("tenant"))
		{
			// Manajemen Menu Tenant
			tenantRoutes.GET("/tenant/menus", menuHandler.GetAll)
			tenantRoutes.POST("/menus", menuHandler.Create)
			tenantRoutes.PUT("/menus/:id", menuHandler.Update)
			tenantRoutes.DELETE("/menus/:id", menuHandler.Delete)

			// Manajemen Pesanan & Laporan Tenant
			tenantRoutes.GET("/tenant/orders", orderHandler.GetTenantOrders)
			tenantRoutes.PUT("/tenant/orders/:id/status", orderHandler.UpdateTenantOrderStatus)
			tenantRoutes.GET("/tenant/earnings", orderHandler.GetTenantEarnings)
		}

		// admin / kasir
		adminRoutes := protected.Group("/")
		adminRoutes.Use(middlewares.RoleMiddleware("admin"))
		{
			adminRoutes.GET("/admin/orders", orderHandler.GetAllOrders)
			adminRoutes.PUT("/admin/orders/:id/pay", orderHandler.MarkOrderAsPaid)

			// Manajemen Tenant oleh Admin
			adminRoutes.GET("/admin/tenants", userHandler.GetTenants)
			adminRoutes.POST("/admin/tenants", userHandler.Register)
		}
	}

	// port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}
