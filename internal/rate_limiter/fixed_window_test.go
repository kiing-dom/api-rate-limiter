package rate_limiter

import (
	"testing"
	"time"
)

func BenchmarkFixedWindow(b *testing.B) {
	rl := NewFixedWindow(10, time.Second*2)
	for i := 0; b.Loop(); i++ {
		rl.Allow()
	}
}
