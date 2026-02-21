package config

import (
	"os"
)

type Config struct {
	Port               string
	RedisURL           string
	UserServiceURL     string
	CatalogServiceURL  string
	CartServiceURL     string
	OrderServiceURL    string
	JWTSecret          string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		RedisURL:          getEnv("REDIS_URL", "localhost:6379"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "localhost:50051"),
		CatalogServiceURL: getEnv("CATALOG_SERVICE_URL", "localhost:50052"),
		CartServiceURL:    getEnv("CART_SERVICE_URL", "localhost:50053"),
		OrderServiceURL:   getEnv("ORDER_SERVICE_URL", "localhost:50055"),
		JWTSecret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
