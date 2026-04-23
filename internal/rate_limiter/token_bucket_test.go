package rate_limiter

import "testing"

func BenchmarkTokenBucket(b *testing.B) {
	rl := NewTokenBucket(10, 1)
	for i := 0; b.Loop(); i++ {
		rl.Allow()
	}
}
