package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	LogLevel        string
	GinPort         int
	EchoPort        int
	DynamoEndpoint  string
	DynamoTable     string
	AWSRegion       string
	JWTSecret       string
	OTELEndpoint    string
	RateLimitBurst  int
	RateLimitPerSec float64
}

func Load() (Config, error) {
	cfg := Config{
		LogLevel:        getEnv("AEVUM_LOG_LEVEL", "info"),
		GinPort:         getEnvInt("AEVUM_GIN_PORT", 8080),
		EchoPort:        getEnvInt("AEVUM_ECHO_PORT", 9090),
		DynamoEndpoint:  os.Getenv("AEVUM_DYNAMODB_ENDPOINT"),
		DynamoTable:     getEnv("AEVUM_DYNAMODB_TABLE", "aevum-events"),
		AWSRegion:       getEnv("AEVUM_AWS_REGION", "eu-central-1"),
		JWTSecret:       os.Getenv("AEVUM_JWT_SECRET"),
		OTELEndpoint:    getEnv("AEVUM_OTEL_ENDPOINT", "localhost:4317"),
		RateLimitBurst:  getEnvInt("AEVUM_RATE_LIMIT_BURST", 100),
		RateLimitPerSec: float64(getEnvInt("AEVUM_RATE_LIMIT_RATE", 50)),
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("missing required env var AEVUM_JWT_SECRET")
	}
	if cfg.GinPort <= 0 || cfg.EchoPort <= 0 {
		return Config{}, fmt.Errorf("invalid ports configured")
	}
	if cfg.RateLimitBurst <= 0 || cfg.RateLimitPerSec <= 0 {
		return Config{}, fmt.Errorf("rate limit values must be greater than zero")
	}
	if cfg.DynamoTable == "" {
		return Config{}, fmt.Errorf("dynamodb table must not be empty")
	}
	if cfg.OTELEndpoint == "" {
		return Config{}, fmt.Errorf("otel endpoint must not be empty")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}
