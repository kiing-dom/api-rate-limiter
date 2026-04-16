package rate_limiter

import (
	"log"
	"sync"
)

var mu sync.Mutex
var Limiters = make(map[string]*RateLimiter)

func GetRateLimiter(userID string, maxTokens float64, refillRate float64) *RateLimiter {
	mu.Lock()
	if rl, exists := Limiters[userID]; exists {
		mu.Unlock()
		return rl
	}

	log.Printf("Creating new rate limiter for user: %s", userID)
	rl := NewRateLimiter(maxTokens, refillRate)
	Limiters[userID] = rl
	mu.Unlock()
	return rl
}
