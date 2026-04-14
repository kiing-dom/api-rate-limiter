package rate_limiter

var Limiters = make(map[string]*RateLimiter)

func GetRateLimiter(userID string, maxTokens float64, refillRate float64) *RateLimiter {
	if rl, exists := Limiters[userID]; exists {
		return rl
	}

	rl := NewRateLimiter(maxTokens, refillRate)
	Limiters[userID] = rl
	return rl
}
