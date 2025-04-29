package config

import (
	"fmt"
	"os"
	"strconv"
)

// AppConfig holds all application configurations
type AppConfig struct {
	Server      ServerConfig
	RateLimiter RateLimiterConfig
	WorkerPool  WorkerPoolConfig
}

// ServerConfig holds HTTP server related configurations
type ServerConfig struct {
	Port string
}

// RateLimiterConfig holds rate limiter related configurations
type RateLimiterConfig struct {
	Enabled       bool
	MaxReqsPerMin int
}

// WorkerPoolConfig holds worker pool related configurations
type WorkerPoolConfig struct {
	WorkerCount int
}

// LoadConfig loads application configurations from environment variables
func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{
		Server: ServerConfig{
			Port: getEnvWithDefault("SERVER_PORT", "8080"),
		},
		RateLimiter: RateLimiterConfig{},
		WorkerPool:  WorkerPoolConfig{},
	}

	// Configure worker pool
	workerCount := getEnvWithDefault("WORKER_COUNT", "500")
	if parsed, err := strconv.Atoi(workerCount); err == nil && parsed > 0 {
		config.WorkerPool.WorkerCount = parsed
	}

	rateLimitStatus := getEnvWithDefault("RATE_LIMITER", "disabled")
	if rateLimitStatus == "enabled" {
		config.RateLimiter.Enabled = true
	} else {
		config.RateLimiter.Enabled = false
	}

	maxRequests := getEnvWithDefault("MAX_REQUESTS_PER_MIN", "10")
	if parsed, err := strconv.Atoi(maxRequests); err == nil && parsed > 0 {
		config.RateLimiter.MaxReqsPerMin = parsed
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// validateConfig checks if all required configurations are set
func validateConfig(config *AppConfig) error {
	if config.RateLimiter.Enabled && config.RateLimiter.MaxReqsPerMin == 0 {
		return fmt.Errorf("rate limiter max requests must be greater than zero")
	}

	if config.WorkerPool.WorkerCount <= 0 {
		return fmt.Errorf("worker count must be greater than zero")
	}

	return nil
}
