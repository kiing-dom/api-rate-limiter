package rate_limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type FixedWindow struct {
	Client *redis.Client
	Limit  int
	Window time.Duration
	UserID string
	Now    func() time.Time
}

func NewFixedWindow(client *redis.Client, userID string, limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		Client: client,
		UserID: userID,
		Limit:  limit,
		Window: window,
		Now:    time.Now,
	}
}

func (rl *FixedWindow) Allow(userID string) bool {
	ctx := context.Background()
	windowSlot := rl.Now().Truncate(rl.Window).Unix()
	key := fmt.Sprintf("ratelimit:fixed:%s:%d", userID, windowSlot)

	count, err := rl.Client.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		rl.Client.Expire(ctx, key, rl.Window)
	}

	return count <= int64(rl.Limit)
}
