package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	EnvPort         = "PORT"
	EnvProductURL   = "PRODUCT_URL"
	EnvProductToken = "PRODUCT_TOKEN"
	EnvRetryCount   = "RETRY_COUNT"
	EnvRetryDelay   = "RETRY_DELAY"
)

type Config struct {
	Server  ServerConfig
	Product ProductConfig
	Retry   RetryConfig
}
type ServerConfig struct {
	Port string
}
type ProductConfig struct {
	Url   string
	Token string
}
type RetryConfig struct {
	Count int
	Delay time.Duration
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv(EnvPort, "8082"),
		},
		Product: ProductConfig{
			Url:   getEnv(EnvProductURL, "http://route256.pavl.uk:8080"),
			Token: getEnv(EnvProductToken, "testtoken"),
		},
		Retry: RetryConfig{
			Count: getEnvInt(EnvRetryCount, 3),
			Delay: getEnvDelay(EnvRetryDelay, 200),
		},
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
		log.Printf("WARN: invalid int value for %s: %s. Using default %d", key, value, fallback)
	}
	return fallback
}

func getEnvDelay(key string, fallback int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return time.Duration(v) * time.Millisecond
		}
		log.Printf("WARN: invalid int value for %s: %s. Using default %d", key, value, fallback)
	}
	return time.Duration(fallback) * time.Millisecond
}
