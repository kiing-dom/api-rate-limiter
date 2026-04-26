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

func NewStore(newAddr string) (*Store, error) {
	client := NewRedisClient(newAddr)
	if err := Ping(context.Background(), client); err != nil {
		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	return &Store{client: client}, nil
}

func (s *Store) GetRateLimiter(userID string) rate_limiter.RateLimiter {
	return rate_limiter.NewTokenBucket(s.client, 3, time.Minute.Seconds()*2)
}
