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