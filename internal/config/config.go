package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	AppName    string
	LogLevel   string
	ServerName string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {
	// Загружаем .env если есть
	_ = godotenv.Load(".env")

	return &Config{
		Port:       getEnv("PORT", ":3000"),
		AppName:    getEnv("APP_NAME", "prtask2"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		ServerName: getEnv("SERVER_NAME", "prtask2"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "prtask2"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
