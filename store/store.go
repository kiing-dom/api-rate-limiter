package store

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/kiing-dom/api-rate-limiter/internal/config"
	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	cfg    *config.Config
}

type KeyConfig struct {
	Algo       string  `json:"algo"`
	RateLimit  int     `json:"rate_limit"`
	WindowSecs int     `json:"window_secs"`
	MaxTokens  float64 `json:"max_tokens"`
	RefillRate float64 `json:"refill_rate"`
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
		if algo != s.cfg.DefaultAlgo {
			return s.GetRateLimiter(s.cfg.DefaultAlgo)
		}

		return rate_limiter.NewTokenBucket(s.client, s.cfg.MaxTokens, s.cfg.RefillRate)
	}
}

func (s *Store) SetKeyConfig(userID string, cfg KeyConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode key config: %w", err)
	}
	return s.client.Set(context.Background(), "keyconfig:"+userID, data, 0).Err()
}

func (s *Store) GetKeyConfig(userID string) (*KeyConfig, error) {
	data, err := s.client.Get(context.Background(), "keyconfig:"+userID).Bytes()
	if err != nil {
		return nil, err
	}

	var cfg KeyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode key config: %w", err)
	}

	return &cfg, nil
}
