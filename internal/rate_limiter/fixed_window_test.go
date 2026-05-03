package rate_limiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func BenchmarkFixedWindow(b *testing.B) {
	db, _ := redismock.NewClientMock()
	userID := "user:abc123"
	rl := NewFixedWindow(db, 10, time.Second*2)
	for i := 0; b.Loop(); i++ {
		rl.Allow(userID)
	}
}

func TestFixedWindow_AllowsWithinLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()
	fixedTime := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
	window := 2 * time.Minute
	userID := "user:abc123"
	windowSlot := fixedTime.Truncate(window).Unix()
	key := fmt.Sprintf("ratelimit:fixed:%s:%d", userID, windowSlot)

	rl := NewFixedWindow(db, 4, window)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, window).SetVal(true) // only happen when the count is == 1

	if !rl.Allow(userID) {
		t.Fatal("Expected request to be allowed")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFixedWindow_RejectsOverLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()
	fixedTime := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	window := 2 * time.Minute
	userID := "user:abc123"
	windowSlot := fixedTime.Truncate(window).Unix()
	key := fmt.Sprintf("ratelimit:fixed:%s:%d", userID, windowSlot)

	limit := 3
	rl := NewFixedWindow(db, limit, window)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectIncr(key).SetVal(int64(limit + 1))

	if rl.Allow(userID) {
		t.Fatal("Expected request to be rejected because over limit")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
