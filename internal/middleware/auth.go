package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"plan-balance-service/pkg/utils"
)

// AuthMiddleware memverifikasi JWT token dari header Authorization.
// Kalau valid, user_id akan disimpan di gin.Context supaya bisa dipakai handler.
// Kalau tidak valid, request langsung ditolak dengan 401.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// 2. Pastikan formatnya "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Authorization header must be: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Verifikasi token pakai fungsi yang sudah ada di pkg/utils/jwt.go
		claims, err := utils.VerifyToken(tokenString, []byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Token is invalid or expired",
				},
			})
			c.Abort()
			return
		}

		// 4. Simpan user_id di context supaya handler bisa akses
		c.Set("user_id", claims.UserID)

		// 5. Lanjut ke handler berikutnya
		c.Next()
	}
}
