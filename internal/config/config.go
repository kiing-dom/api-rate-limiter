package config

import (
	"time"
)

type Config struct {
	RedisAddr  string
	HTTPPort   string
	GRPCPort   string
	RateLimit  int
	Window     time.Duration
	RefillRate float64
}
