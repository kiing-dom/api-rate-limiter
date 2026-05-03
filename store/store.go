package store

import (
	"context"
	"fmt"
	"time"

	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

type RLStore interface {
	GetRateLimiter(algo string) rate_limiter.RateLimiter
}

func NewStore(newAddr string) (*Store, error) {
	client := NewRedisClient(newAddr)
	if err := Ping(context.Background(), client); err != nil {
		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	return &Store{client: client}, nil
}

func (s *Store) GetRateLimiter(algo string) rate_limiter.RateLimiter {
	switch algo {
	case "sliding":
		return rate_limiter.NewSlidingWindow(s.client, 3, 2*time.Minute)
	case "fixed":
		return rate_limiter.NewFixedWindow(s.client, 3, 2*time.Minute)
	case "token":
		return rate_limiter.NewTokenBucket(s.client, 3, 1)
	default:
		return rate_limiter.NewTokenBucket(s.client, 3, 1)
	}
}
