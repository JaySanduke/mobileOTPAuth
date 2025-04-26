package api

import (
	"mobileOTPAuth/config"
	"mobileOTPAuth/internal/model"
	"mobileOTPAuth/internal/redis"
	"mobileOTPAuth/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redisPkg "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LoginRequest struct {
	MobileNumber string `json:"mobile_number" binding:"required"`
	OTP          string `json:"otp" binding:"required"`
	Fingerprint  string `json:"fingerprint" binding:"required"`
}

func LoginHandler(db *gorm.DB, redisClient *redis.RedisClient, jwtManagerClient *utils.JWTManager, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		if loginRequest.MobileNumber == "" || loginRequest.OTP == "" || loginRequest.Fingerprint == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		storedOTP, err := redisClient.GetOTP(loginRequest.MobileNumber)
		if err == redisPkg.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not sent"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve OTP"})
			return
		}

		if loginRequest.OTP != storedOTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		}

		var user model.User
		result := db.Model(&model.User{}).Where("mobile_number = ?", loginRequest.MobileNumber).First(&user)
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if user.Fingerprint != loginRequest.Fingerprint {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Fingerprint mismatch"})
			return
		}

		token, err := jwtManagerClient.GenerateToken(user.ID, user.MobileNumber, loginRequest.Fingerprint, 2*60)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		redisClient.DeleteOTP(loginRequest.MobileNumber)

		userLoginSession := model.UserLoginSession{
			UserID:       strconv.FormatUint(uint64(user.ID), 10),
			MobileNumber: loginRequest.MobileNumber,
			Fingerprint:  loginRequest.Fingerprint,
			LastLogin:    time.Now(),
			AccessToken:  token,
		}
		if err := db.Model(&model.UserLoginSession{}).Create(&userLoginSession).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create login session in the db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"message":     "Login successful",
			"accessToken": token,
			"user":        user,
		})
	}
}
