package store

import (
	"context"
	"fmt"
	"time"

	"github.com/kiing-dom/api-rate-limiter/internal/config"
	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	cfg    *config.Config
}

type RLStore interface {
	GetRateLimiter(algo string) rate_limiter.RateLimiter
}

func NewStore(newAddr string, cfg *config.Config) (*Store, error) {
	client := NewRedisClient(newAddr)
	if err := Ping(context.Background(), client); err != nil {
		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	return &Store{client: client, cfg: cfg}, nil
}

func (s *Store) GetRateLimiter(algo string) rate_limiter.RateLimiter {
	switch algo {
	case "sliding":
		return rate_limiter.NewSlidingWindow(s.client, s.cfg.RateLimit, s.cfg.Window)
	case "fixed":
		return rate_limiter.NewFixedWindow(s.client, s.cfg.RateLimit, s.cfg.Window)
	case "token":
		return rate_limiter.NewTokenBucket(s.client, s.cfg.MaxTokens, s.cfg.RefillRate)
	default:
		return rate_limiter.NewTokenBucket(s.client, s.cfg.MaxTokens, s.cfg.RefillRate)
	}
}
