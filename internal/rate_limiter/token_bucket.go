package rate_limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	Client     *redis.Client
	MaxTokens  float64 // max bucket size (burst capacity)
	RefillRate float64 // how many tokens come back per second
}

func NewTokenBucket(client *redis.Client, maxTokens float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		Client:     client,
		MaxTokens:  maxTokens,
		RefillRate: refillRate,
	}
}

func (rl *TokenBucket) Allow(userID string) bool {
	// add context
	ctx := context.Background()
	// create key (using userID)
	key := fmt.Sprintf("ratelimit:token:%s", userID)
	// get time
	now := time.Now()

	// get values from redis (err if no value)
	vals, err := rl.Client.HMGet(ctx, key, "tokens", "last_refill").Result()
	if err != nil {
		return false
	}

	// initialize values for:
	// - tokens
	tokens := rl.MaxTokens
	// - lastRefill
	lastRefill := now

	// if they are in the store update the values to the store values
	if vals[0] != nil && vals[1] != nil {
		tokens, _ = strconv.ParseFloat(vals[0].(string), 64)
		lastRefillNano, _ := strconv.ParseInt(vals[1].(string), 10, 64)
		lastRefill = time.Unix(0, lastRefillNano)
	}

	elapsed := now.Sub(lastRefill).Seconds()
	tokens = min(rl.MaxTokens, tokens+elapsed*rl.RefillRate)

	// if we dont have enough tokens (tokens < 1) return false
	if tokens < 1 {
		return false
	}
	// otherwise it's allowed
	// -> decrement token count by 1
	tokens -= 1
	// update the redis store
	rl.Client.HSet(ctx, key, "tokens", tokens, "last_refill", now.UnixNano())
	ttl := time.Duration(rl.MaxTokens/rl.RefillRate)*time.Second + time.Minute
	rl.Client.Expire(ctx, key, ttl)

	return true
}
