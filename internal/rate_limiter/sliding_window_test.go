package rate_limiter

import (
	"testing"
	"time"
)

func BenchmarkSlidingWindow(b *testing.B) {
	rl := NewSlidingWindow(10, time.Second*2)
	for i := 0; b.Loop(); i++ {
		rl.Allow()
	}
}
