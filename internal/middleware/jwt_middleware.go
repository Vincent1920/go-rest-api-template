package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"todo-api/internal/dto"
	"todo-api/internal/utils"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Ambil Authorization Header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse("Authorization header is required", nil))
			return
		}

		// Format harus: Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse("Invalid authorization format", nil))
			return
		}

		// Validasi JWT
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.ErrorResponse("Invalid or expired token", nil))
			return
		}

		// Simpan userID ke context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
