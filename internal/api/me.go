package api

import (
	"mobileOTPAuth/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user model.User
		result := db.Model(&model.User{}).Where("id = ?", userID).First(&user)
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"id":          user.ID,
			"mobile":      user.MobileNumber,
			"fingerprint": user.Fingerprint,
			"created_at":  user.CreatedAt,
		})
	}
}
