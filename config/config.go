package config

import (
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	RedisPass  string
	JWTSecret  string
}

func LoadConfig() *Config {
	// Uncomment the godotenv load method in the development mode
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("No .env file found, using system environment variables")
	// }

	return &Config{
		Port:       os.Getenv("PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RedisAddr:  os.Getenv("REDIS_ADDR"),
		RedisPass:  os.Getenv("REDIS_PASSWORD"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}
