package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DatabasePath   string
	MaxFileSize    int64
	HealthCheckURL string
	LogLevel       string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "3000"),
		DatabasePath:   getEnv("DATABASE_PATH", "./proxies.db"),
		MaxFileSize:    getEnvInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB default
		HealthCheckURL: getEnv("HEALTH_CHECK_URL", "https://httpbin.org/ip"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
