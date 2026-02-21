package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	SMTPHost    string
	SMTPPort    string
	SMTPFrom    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "50057"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/notification_db?sslmode=disable"),
		SMTPHost:    getEnv("SMTP_HOST", "localhost"),
		SMTPPort:    getEnv("SMTP_PORT", "1025"),
		SMTPFrom:    getEnv("SMTP_FROM", "noreply@microservices-demo.local"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
