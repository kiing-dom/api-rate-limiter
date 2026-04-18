package rate_limiter

import (
	"sync"
	"time"
)

type FixedWindow struct {
	mu          sync.Mutex
	Count       int
	Limit       int
	Window      time.Duration
	WindowStart time.Time
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		Limit:       limit,
		Window:      window,
		WindowStart: time.Now(),
	}
}

func (rl *FixedWindow) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if now.Sub(rl.WindowStart) > rl.Window {
		rl.WindowStart = now
		rl.Count = 0
	}

	if rl.Count < rl.Limit {
		rl.Count++
		return true
	}

	return false
}
