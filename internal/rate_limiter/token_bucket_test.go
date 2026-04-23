package rate_limiter

import (
	"testing"
	"time"
)

func BenchmarkTokenBucket(b *testing.B) {
	rl := NewTokenBucket(10, 1)
	for i := 0; b.Loop(); i++ {
		rl.Allow()
	}
}

func TestTokenBucket_AllowsWithinLimit(t *testing.T) {
	rl := NewTokenBucket(5, 1)

	for i := range 5 {
		if !rl.Allow() {
			t.Fatalf("Expected request %d to be allowed", i)
		}
	}
}

func TestTokenBucket_RejectsOverLimit(t *testing.T) {
	rl := NewTokenBucket(5, 1)

	for range 5 {
		rl.Allow()
	}

	if rl.Allow() {
		t.Fatal("Expected request to be rejected!")
	}
}

func TestTokenBucket_RefillsOverTime(t *testing.T) {
	rl := NewTokenBucket(1, 1)

	if !rl.Allow() {
		t.Fatal("Expected first request to be allowed")
	}

	if rl.Allow() {
		t.Fatal("Expected second request to be rejected")
	}

	time.Sleep(time.Second)

	if !rl.Allow() {
		t.Fatal("Expected third request to be allowed (following refill)")
	}
}
