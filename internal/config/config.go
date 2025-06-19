package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	LogLevel    string
	JWTSecret   string
	
	// OpenTelemetry configuration
	OtelEnabled      bool
	OtelServiceName  string
	OtelExporterType string
	OtelEndpoint     string
	
	// Temporal configuration
	TemporalHost      string
	TemporalNamespace string
	TaskQueue         string
	
	// OPA configuration
	OPAEnabled bool
	OPAURL     string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://yugabyte:yugabyte@localhost:5433/yugabyte?sslmode=disable"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		
		// OpenTelemetry configuration
		OtelEnabled:      getEnvBool("OTEL_ENABLED", false),
		OtelServiceName:  getEnv("OTEL_SERVICE_NAME", "golang-boilerplate"),
		OtelExporterType: getEnv("OTEL_EXPORTER_TYPE", "jaeger"),
		OtelEndpoint:     getEnv("OTEL_ENDPOINT", "http://localhost:14268/api/traces"),
		
		// Temporal configuration
		TemporalHost:      getEnv("TEMPORAL_HOST", ""),
		TemporalNamespace: getEnv("TEMPORAL_NAMESPACE", "default"),
		TaskQueue:         getEnv("TASK_QUEUE", "user-onboarding"),
		
		// OPA configuration
		OPAEnabled: getEnvBool("OPA_ENABLED", false),
		OPAURL:     getEnv("OPA_URL", "http://localhost:8181"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}