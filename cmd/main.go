package main

import (
	"fmt"
	"time"

	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
)

func main() {
	userOne := "abc123"
	rlOne := rate_limiter.NewRateLimiter(5, 1)
	rate_limiter.Limiters[userOne] = rlOne

	users := []string{"abc123", "def456"}
	// user_two := "def456"

	// rate_limiter.Limiters[user_two] = nil

	for i := range users {
		fmt.Printf("User %d: %s\n", i, users[i])

		rl := rate_limiter.GetRateLimiter(users[i], 3, 0.33)

		for j := 1; j <= 6; j++ {
			fmt.Printf("Attempt %d. Request Allowed?: %v\n", j, rl.Allow())
			time.Sleep(time.Second)
		}
	}
}
