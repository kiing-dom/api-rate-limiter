package rate_limiter

import "log"

var Limiters = make(map[string]*RateLimiter)

func GetRateLimiter(userID string, maxTokens float64, refillRate float64) *RateLimiter {
	if rl, exists := Limiters[userID]; exists {
		return rl
	}

	log.Printf("Creating new rate limiter for user: %s", userID)
	rl := NewRateLimiter(maxTokens, refillRate)
	Limiters[userID] = rl
	return rl
}
