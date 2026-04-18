package rate_limiter

import "time"

func NewRateLimiter(algo string) RateLimiter {
	switch algo {
	case "token":
		return NewTokenBucket(5, 1)
	case "fixed":
		return NewFixedWindow(5, time.Minute*2)
	// add sliding window when implemented
	default:
		panic("unknown algorithm provided")
	}
}
