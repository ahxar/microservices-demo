package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	RedisURL   string
	CartTTLDays int
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "50053"),
		RedisURL:   getEnv("REDIS_URL", "localhost:6379"),
		CartTTLDays: getEnvAsInt("CART_TTL_DAYS", 7),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
