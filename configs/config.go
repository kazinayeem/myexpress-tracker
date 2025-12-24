package configs

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	ServerPort     string
	DatabasePath   string
	JWTSecret      string
	JWTExpiration  time.Duration
	Environment    string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	jwtExpHours := getEnvAsInt("JWT_EXPIRATION_HOURS", 24)
	
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DatabasePath:   getEnv("DATABASE_PATH", "./data/tracker.db"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration:  time.Duration(jwtExpHours) * time.Hour,
		Environment:    getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer with fallback
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
