package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
)

func TestRedisTokenBucket_AllowsWithinLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()

	userID := "user:abc123"
	key := fmt.Sprintf("ratelimit:token:%s", userID)
	maxTokens := 3

	fixedTime := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	rl := rate_limiter.NewTokenBucket(db, float64(maxTokens), 1)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectHMGet(key, "tokens", "last_refill").SetVal([]interface{}{nil, nil})
	mock.ExpectHSet(key, "tokens", float64(maxTokens-1), "last_refill", fixedTime.UnixNano()).SetVal(1)
	ttl := time.Duration(maxTokens)*time.Second + time.Minute
	mock.ExpectExpire(key, ttl).SetVal(true)

	if !rl.Allow(userID) {
		t.Fatal("Expected first request to be allowed")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
