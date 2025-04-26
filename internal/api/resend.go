package api

import (
	"mobileOTPAuth/config"
	"mobileOTPAuth/internal/model"
	"mobileOTPAuth/internal/redis"
	"mobileOTPAuth/internal/service"
	"mobileOTPAuth/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResendRequest struct {
	MobileNumber string `json:"mobile_number" binding:"required"`
	Fingerprint  string `json:"fingerprint" binding:"required"`
}

func ResendOTPHandler(db *gorm.DB, redisClient *redis.RedisClient, smsServiceClient service.SMSService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req ResendRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user model.User
		result := db.Model(&model.User{}).Where("mobile_number = ?", req.MobileNumber).First(&user)
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User account not found"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		otp := utils.GenerateOTP(6)
		if err := redisClient.SetOTP(req.MobileNumber, otp, time.Duration(2)*time.Minute); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
			return
		}

		// Dummy SMS provider to send OTP via message
		err := smsServiceClient.SendOTP(c, req.MobileNumber, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   "OTP resent successfully",
			"debug_otp": otp, // remove this in the production env
		})
	}
}
