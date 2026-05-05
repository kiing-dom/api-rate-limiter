package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr  string
	HTTPPort   string
	GRPCPort   string
	RateLimit  int
	Window     time.Duration
	MaxTokens  float64
	RefillRate float64
}

func Load() (*Config, error) {
	godotenv.Load("../../.env")

	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "3"))
	windowSecs, _ := strconv.Atoi(getEnv("WINDOW_SECONDS", "120"))
	maxTokens, _ := strconv.ParseFloat(getEnv("MAX_TOKENS", "5"), 64)
	refillRate, _ := strconv.ParseFloat(getEnv("REFILL_RATE", "1.0"), 64)

	return &Config{
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		HTTPPort:   getEnv("HTTP_PORT", "8081"),
		GRPCPort:   getEnv("GRPC_PORT", "50051"),
		RateLimit:  rateLimit,
		Window:     time.Duration(windowSecs) * time.Second,
		MaxTokens:  maxTokens,
		RefillRate: refillRate,
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
