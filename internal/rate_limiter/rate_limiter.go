package rate_limiter

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	Tokens     float64 // how many requests can be made
	MaxTokens  float64 // max bucket size (burst capacity)
	RefillRate float64 // how fast tokens come back
	LastRefill time.Time
}

func NewRateLimiter(tokens float64, refillRate float64) *RateLimiter {
	return &RateLimiter{
		Tokens:     tokens,
		MaxTokens:  5,
		RefillRate: refillRate,
		LastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() string {
	elapsed := time.Since(rl.LastRefill)
	tokensToAdd := elapsed.Seconds() * rl.RefillRate

	rl.Tokens = min(rl.MaxTokens, rl.Tokens+tokensToAdd)

	rl.LastRefill = time.Now()

	if rl.Tokens >= 1 {
		rl.Tokens -= 1
		return fmt.Sprintf("Request allowed! %f requests remaining \n", rl.Tokens)
	} else {
		return "Rejected. No remaining tokens"
	}
}
