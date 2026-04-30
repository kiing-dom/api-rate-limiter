package rate_limiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func BenchmarkTokenBucket(b *testing.B) {
	db, _ := redismock.NewClientMock()
	rl := NewTokenBucket(db, float64(3), float64(1))
	userID := "user:abc123"
	for i := 0; b.Loop(); i++ {
		rl.Allow(userID)
	}
}

func TestRedisTokenBucket_AllowsWithinLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()

	userID := "user:abc123"
	key := fmt.Sprintf("ratelimit:token:%s", userID)
	maxTokens := 3

	fixedTime := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	rl := NewTokenBucket(db, float64(maxTokens), 1)
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

func TestTokenBucket_RejectsOverLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()

	userID := "user:cde456"
	key := fmt.Sprintf("ratelimit:token:%s", userID)
	maxTokens := 3

	fixedTime := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	rl := NewTokenBucket(db, float64(maxTokens), 1)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectHMGet(key, "tokens", "last_refill").SetVal([]interface{}{
		"0",
		fmt.Sprintf("%d", fixedTime.UnixNano()),
	})

	if rl.Allow(userID) {
		t.Fatal("Expected request to be rejected")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTokenBucket_RefillsOverTime(t *testing.T) {
	db, mock := redismock.NewClientMock()

	userID := "user:cde456"
	key := fmt.Sprintf("ratelimit:token:%s", userID)
	maxTokens := 1

	fixedTime := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
	pastTime := fixedTime.Add(2 * time.Second)

	rl := NewTokenBucket(db, float64(maxTokens), 1)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectHMGet(key, "tokens", "last_refill").SetVal([]interface{}{
		"0",
		fmt.Sprintf("%d", pastTime.UnixNano()),
	})

	mock.ExpectHSet(key, "tokens", float64(0), "last_refill", fixedTime.UnixNano()).SetVal(1)
	ttl := time.Duration(maxTokens)*time.Second + time.Minute
	mock.ExpectExpire(key, ttl).SetVal(true)

	if !rl.Allow(userID) {
		t.Fatal("Expected request to be allowed after refill")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
