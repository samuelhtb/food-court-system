package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/services"
)

// UserHandler adalah struct yang menjembatani jalur web (Gin) dengan logika bisnis (Service)
type UserHandler struct {
	userService services.UserService
	secretKey   string
}

// NewUserHandler berfungsi untuk merakit handler ini
func NewUserHandler(userService services.UserService, secretKey string) *UserHandler {
	return &UserHandler{userService, secretKey}
}

// Register menerima permintaan pendaftaran akun baru
func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	
	// 1. Tangkap data JSON dari request dan masukkan ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// 2. Suruh Service untuk menyimpan data tersebut
	if err := h.userService.RegisterTenant(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Kembalikan respon sukses
	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil, silakan login"})
}

// Login menerima email & password, lalu mengembalikan Token JWT
func (h *UserHandler) Login(c *gin.Context) {
	// Struct sementara untuk menangkap input login
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 1. Tangkap input login
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password wajib diisi"})
		return
	}

	// 2. Suruh Service untuk mengecek ke database
	user, err := h.userService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 3. Jika berhasil, buatkan Token (Tiket Masuk) JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	// 4. Stempel token tersebut dengan Kunci Rahasia
	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token otorisasi"})
		return
	}

	// 5. Berikan token ke pengguna
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"role":    user.Role,
		"name":    user.Name,
	})
}