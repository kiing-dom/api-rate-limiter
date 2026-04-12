package cmd

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens     float64 // how many requests can be made
	maxTokens  float64 // max bucket size (burst capacity)
	refillRate float64 // how fast tokens come back
	lastRefill time.Time
}

func (rl *RateLimiter) Allow() {
	elapsed := time.Since(rl.lastRefill)
	tokensToAdd := elapsed.Seconds() * rl.refillRate
	tokens := rl.tokens

	tokens = min(rl.maxTokens, tokens+tokensToAdd)

	rl.lastRefill = time.Now()

	if tokens >= 1 {
		tokens -= 1
		fmt.Printf("Request allowed! %f requests remaining", tokens)
	} else {
		fmt.Println("Rejected. No remaining tokens")
	}
}
