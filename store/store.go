package store

import (
	"sync"
	"time"

	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
)

type Store struct {
	mu           sync.Mutex
	RateLimiters map[string]rate_limiter.RateLimiter
}

func NewStore() *Store {
	return &Store{
		RateLimiters: make(map[string]rate_limiter.RateLimiter),
	}
}

func (s *Store) GetRateLimiter(userID string) rate_limiter.RateLimiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	rl, exists := s.RateLimiters[userID]
	if !exists {
		rl := rate_limiter.NewFixedWindow(5, time.Minute*2)
		s.RateLimiters[userID] = rl
	}

	return rl
}
