package rate_limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SlidingWindow struct {
	Client *redis.Client
	Limit  int
	Window time.Duration
	Now    func() time.Time
}

func NewSlidingWindow(client *redis.Client, limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		Client: client,
		Limit:  limit,
		Window: window,
		Now:    time.Now,
	}
}

func (rl *SlidingWindow) Allow(userID string) bool {
	ctx := context.Background()
	now := rl.Now()
	key := fmt.Sprintf("ratelimit:sliding:%s", userID)
	cutoff := fmt.Sprintf("%d", now.Add(-rl.Window).UnixNano())

	pipe := rl.Client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", cutoff)
	countCmd := pipe.ZCard(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil {
		return false
	}

	if countCmd.Val() >= int64(rl.Limit) {
		return false
	}

	nano := now.Unix()
	rl.Client.ZAdd(ctx, key, redis.Z{Score: float64(nano), Member: nano})
	rl.Client.Expire(ctx, key, rl.Window)
	return true
}
