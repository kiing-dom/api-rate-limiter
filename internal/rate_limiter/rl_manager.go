package rate_limiter

var limiters = make(map[string]*RateLimiter)

func GetRateLimiter(userID string, maxTokens float64, refillRate float64) *RateLimiter {
	if rl, exists := limiters[userID]; exists {
		return rl
	}

	rl := NewRateLimiter(maxTokens, refillRate)
	limiters[userID] = rl
	return rl
}
