// Package config provides configuration management for the Jira Clone Backend application.
// It handles loading and parsing environment variables into a structured configuration object.
package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application.
// It contains server settings, JWT configuration, logging levels, and rate limiting options.
type Config struct {
	Port      string // Server port (default: "8080")
	LogLevel  string // Logging level: debug, info, warn, error, fatal (default: "info")
	JWTSecret string // JWT signing secret for token generation and validation
	RateLimit int    // Rate limit per minute for API requests (default: 100)
}

// Load loads configuration from environment variables and returns a Config struct.
// If environment variables are not set, it uses sensible defaults.
//
// Environment variables:
//   - PORT: Server port (default: "8080")
//   - LOG_LEVEL: Logging level (default: "info")
//   - JWT_SECRET: JWT signing secret (default: "your-secret-key")
//   - RATE_LIMIT: Rate limit per minute (default: 100)
//
// Returns a pointer to the Config struct with loaded values.
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		RateLimit: getEnvAsInt("RATE_LIMIT", 100),
	}
}

// getEnv gets an environment variable with a default value.
// If the environment variable is not set or is empty, it returns the default value.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: The default value to return if the environment variable is not set
//
// Returns the environment variable value or the default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value.
// If the environment variable is not set, is empty, or cannot be parsed as an integer,
// it returns the default value.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: The default integer value to return if parsing fails
//
// Returns the parsed integer value or the default value.
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
