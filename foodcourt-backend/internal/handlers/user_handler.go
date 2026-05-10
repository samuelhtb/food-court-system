package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
	secretKey   string
}

func NewUserHandler(userService services.UserService, secretKey string) *UserHandler {
	return &UserHandler{userService, secretKey}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest // Get data from request body into DTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid atau kurang lengkap"})
		return
	}

	if err := h.userService.RegisterTenant(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil, silakan login"})
}

func (h *UserHandler) RegisterAdmin(c *gin.Context) {
	// Jika sudah login sebagai admin (lewat middleware), bypass pengecekan secret key header
	userRole, _ := c.Get("role")
	secret := c.GetHeader("X-Admin-Secret")

	if userRole != "admin" && secret != os.Getenv("ADMIN_SECRET_KEY") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses dilarang! Memerlukan role admin atau secret key valid."})
		return
	}

	var req dto.RegisterRequest // Get data from request body into DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid atau kurang lengkap"})
		return
	}

	if err := h.userService.RegisterAdmin(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin berhasil dibuat"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email dan password wajib diisi"})
		return
	}

	user, err := h.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 168).Unix(), // Token berlaku 7 hari
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token otorisasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"role":    user.Role,
		"name":    user.Name,
	})
}

func (h *UserHandler) GetTenants(c *gin.Context) {
	tenants, err := h.userService.GetTenants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.UserResponse
	for _, t := range tenants {
		response = append(response, dto.UserResponse{
			ID:         t.ID,
			Username:   t.Username,
			Email:      t.Email,
			Name:       t.Name,
			Role:       t.Role,
			TenantName: t.TenantName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *UserHandler) GetAdmins(c *gin.Context) {
	admins, err := h.userService.GetAdmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.UserResponse
	for _, a := range admins {
		response = append(response, dto.UserResponse{
			ID:       a.ID,
			Username: a.Username,
			Email:    a.Email,
			Name:     a.Name,
			Role:     a.Role,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}