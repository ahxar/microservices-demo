package config

import (
	"os"
)

type Config struct {
	Port                    string
	DatabaseURL             string
	CatalogServiceURL       string
	CartServiceURL          string
	PaymentServiceURL       string
	ShippingServiceURL      string
	NotificationServiceURL  string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "50055"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/order_db?sslmode=disable"),
		CatalogServiceURL:      getEnv("CATALOG_SERVICE_URL", "localhost:50052"),
		CartServiceURL:         getEnv("CART_SERVICE_URL", "localhost:50053"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "localhost:50056"),
		ShippingServiceURL:     getEnv("SHIPPING_SERVICE_URL", "localhost:50058"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "localhost:50057"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
