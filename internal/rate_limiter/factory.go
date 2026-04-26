package rate_limiter

import (
	"github.com/redis/go-redis/v9"
)

func NewRateLimiter(algo string, client *redis.Client) RateLimiter {
	switch algo {
	case "token":
		return NewTokenBucket(client, 5, 1)
	case "fixed":
		// return NewFixedWindow(5, time.Minute*2)
	// add sliding window when implemented
	default:
		panic("unknown algorithm provided")
	}
	return nil
}
