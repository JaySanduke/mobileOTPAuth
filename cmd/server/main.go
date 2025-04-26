package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"mobileOTPAuth/config"
	"mobileOTPAuth/internal/api"
	"mobileOTPAuth/internal/db"
	"mobileOTPAuth/internal/middleware"
	"mobileOTPAuth/internal/redis"
	"mobileOTPAuth/internal/service"
	"mobileOTPAuth/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type serverConfiguration struct {
	PgClient         *gorm.DB
	RedisClient      *redis.RedisClient
	JwtManagerClient *utils.JWTManager
	SmsServiceClient service.SMSService
}

func initializeClients(cfg *config.Config) *serverConfiguration {

	postgreDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	pgClient, err := db.InitPostgres(postgreDSN)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	redisClient := redis.NewRedis(cfg.RedisAddr)

	jwtManagerClient := utils.NewJWTManager(cfg.JWTSecret)

	smsServiceClient := service.NewSMSService()

	return &serverConfiguration{
		PgClient:         pgClient,
		RedisClient:      redisClient,
		JwtManagerClient: jwtManagerClient,
		SmsServiceClient: smsServiceClient,
	}

}

func main() {

	cfg := config.LoadConfig()

	clients := initializeClients(cfg)

	ginRouter := gin.Default()

	public := ginRouter.Group("/api/auth")
	{
		public.POST("/register", api.RegisterHandler(clients.PgClient, clients.RedisClient, clients.SmsServiceClient, cfg))
		public.POST("/login", api.LoginHandler(clients.PgClient, clients.RedisClient, clients.JwtManagerClient, cfg))
		public.POST("/resend-otp", api.ResendOTPHandler(clients.PgClient, clients.RedisClient, clients.SmsServiceClient, cfg))
	}

	protected := ginRouter.Group("/api/auth")
	protected.Use(middleware.JWTAuth(clients.JwtManagerClient))
	{
		protected.GET("/me", api.MeHandler(clients.PgClient))
	}

	ginRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	port := cfg.Port
	if cfg.Port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := ginRouter.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
