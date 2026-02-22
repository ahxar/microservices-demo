package config

import (
	"os"
)

type Config struct {
	Port                     string
	MetricsPort              string
	RedisURL                 string
	UserServiceURL           string
	CatalogServiceURL        string
	CartServiceURL           string
	OrderServiceURL          string
	JWTSecret                string
	OTELExporterOTLPEndpoint string
	OTELExporterOTLPInsecure bool
	OTELServiceName          string
}

func Load() *Config {
	return &Config{
		Port:                     getEnv("PORT", "8080"),
		MetricsPort:              getEnv("METRICS_PORT", "9090"),
		RedisURL:                 getEnv("REDIS_URL", "localhost:6379"),
		UserServiceURL:           getEnv("USER_SERVICE_URL", "localhost:50051"),
		CatalogServiceURL:        getEnv("CATALOG_SERVICE_URL", "localhost:50052"),
		CartServiceURL:           getEnv("CART_SERVICE_URL", "localhost:50053"),
		OrderServiceURL:          getEnv("ORDER_SERVICE_URL", "localhost:50055"),
		JWTSecret:                getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		OTELExporterOTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "tempo:4317"),
		OTELExporterOTLPInsecure: getEnvAsBool("OTEL_EXPORTER_OTLP_INSECURE", true),
		OTELServiceName:          getEnv("OTEL_SERVICE_NAME", "gateway"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "1" || value == "true" || value == "TRUE" || value == "True" {
			return true
		}
		if value == "0" || value == "false" || value == "FALSE" || value == "False" {
			return false
		}
	}

	return defaultValue
}
