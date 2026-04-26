package rate_limiter

type RateLimiter interface {
	Allow(userID string) bool
}
