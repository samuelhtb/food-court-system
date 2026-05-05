package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, butuh token!"})
			c.Abort()
			return
		}

		// Peningkatan: TrimPrefix lebih aman dan rapi daripada Replace
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kedaluwarsa"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// PERBAIKAN KRUSIAL: Gunakan "user_id" agar sinkron dengan menu_handler.go
			c.Set("user_id", claims["user_id"]) 
			
			// Ubah juga menjadi "role" agar seragam format penamaannya (snake_case)
			c.Set("role", claims["role"]) 
		}

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role yang tadi disimpan oleh AuthMiddleware
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Role tidak teridentifikasi"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		isAllowed := false

		// Cek apakah role user ada di dalam daftar role yang diizinkan
		for _, role := range allowedRoles {
			if roleStr == role {
				isAllowed = true
				break
			}
		}

		// Jika tidak cocok, tendang keluar
		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Anda tidak memiliki izin untuk halaman ini"})
			c.Abort()
			return
		}

		c.Next()
	}
}