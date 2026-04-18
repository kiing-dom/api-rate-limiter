package rate_limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	Tokens     float64 // how many requests can be made
	MaxTokens  float64 // max bucket size (burst capacity)
	RefillRate float64 // how many tokens come back per second
	LastRefill time.Time
}

func NewTokenBucket(maxTokens float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		Tokens:     maxTokens,
		MaxTokens:  maxTokens,
		RefillRate: refillRate,
		LastRefill: time.Now(),
	}
}

func (rl *TokenBucket) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
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
