package rate_limiter

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	mu         sync.Mutex
	timestamps []time.Time
	window     time.Duration
	limit      int
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		timestamps: []time.Time{},
		window:     window,
		limit:      limit,
	}
}

func (rl *SlidingWindow) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// remove the old timestamps
	var valid []time.Time
	for _, t := range rl.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	rl.timestamps = valid

	// check limit
	if len(rl.timestamps) > rl.limit {
		return false
	}

	// do current request
	rl.timestamps = append(rl.timestamps, now)

	return true
}
