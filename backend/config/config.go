package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using default environment variables")
	}

	AppConfig = &Config{
		JWTSecret:  getEnv("JWT_SECRET", "mysecretkey123"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "gotasker"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
