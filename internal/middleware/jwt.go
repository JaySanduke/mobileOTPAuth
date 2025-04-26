package middleware

import (
	"errors"
	"mobileOTPAuth/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(jwtManagerClient *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtManagerClient.VerifyToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenInvalidClaims) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("fingerprint", claims.Fingerprint)

		c.Next()
	}
}
