package rate_limiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func BenchmarkSlidingWindow(b *testing.B) {
	db, _ := redismock.NewClientMock()
	userID := "user:abc123"
	rl := NewSlidingWindow(db, userID, 10, time.Second*2)
	for i := 0; b.Loop(); i++ {
		rl.Allow(userID)
	}
}

func TestSlidingWindow_AllowsWithinLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()
	fixedTime := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	window := 2 * time.Minute
	userID := "user:abc123"
	key := fmt.Sprintf("ratelimit:sliding:%s", userID)
	cutoff := fmt.Sprintf("%d", fixedTime.Add(-window).UnixNano())
	nano := fixedTime.Unix()

	limit := 3
	rl := NewSlidingWindow(db, userID, limit, window)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectTxPipeline()
	mock.ExpectZRemRangeByScore(key, "0", cutoff).SetVal(0)
	mock.ExpectZCard(key).SetVal(0)
	mock.ExpectTxPipelineExec()
	mock.ExpectZAdd(key, redis.Z{Score: float64(nano), Member: nano}).SetVal(1)
	mock.ExpectExpire(key, window).SetVal(true)

	if !rl.Allow(userID) {
		t.Fatal("Expected request to be accepted as it is within limit")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSlidingWindow_RejectsOverLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()
	fixedTime := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	window := 2 * time.Minute
	userID := "user:abc123"
	key := fmt.Sprintf("ratelimit:sliding:%s", userID)
	cutoff := fmt.Sprintf("%d", fixedTime.Add(-window).UnixNano())

	limit := 3
	rl := NewSlidingWindow(db, userID, limit, window)
	rl.Now = func() time.Time { return fixedTime }

	mock.ExpectTxPipeline()
	mock.ExpectZRemRangeByScore(key, "0", cutoff).SetVal(0)
	mock.ExpectZCard(key).SetVal(int64(limit))
	mock.ExpectTxPipelineExec()

	if rl.Allow(userID) {
		t.Fatal("Expected request to be rejected as it's over the limit")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
