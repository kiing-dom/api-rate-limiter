package rate_limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	Tokens     float64 // how many requests can be made
	MaxTokens  float64 // max bucket size (burst capacity)
	RefillRate float64 // how many tokens come back per second
	LastRefill time.Time
}

func NewRateLimiter(maxTokens float64, refillRate float64) *RateLimiter {
	return &RateLimiter{
		Tokens:     maxTokens,
		MaxTokens:  maxTokens,
		RefillRate: refillRate,
		LastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	mu.Lock()
	defer mu.Unlock()
	elapsed := time.Since(rl.LastRefill)
	tokensToAdd := elapsed.Seconds() * rl.RefillRate

	rl.Tokens = min(rl.MaxTokens, rl.Tokens+tokensToAdd)
	rl.LastRefill = time.Now()

	if rl.Tokens >= 1 {
		rl.Tokens -= 1
		return true
	} else {
		return false
	}
}
