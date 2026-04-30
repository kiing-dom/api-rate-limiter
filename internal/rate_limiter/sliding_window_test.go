package rate_limiter

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func BenchmarkSlidingWindow(b *testing.B) {
	db, _ := redismock.NewClientMock()
	userID := "user:abc123"
	rl := NewSlidingWindow(db, userID, 10, time.Second*2)
	for i := 0; b.Loop(); i++ {
		rl.Allow(userID)
	}
}
