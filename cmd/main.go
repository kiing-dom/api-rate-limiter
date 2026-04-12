package main

import (
	"fmt"

	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
)

func main() {
	rl := rate_limiter.NewRateLimiter(5, 1)

	for i := 0; i < 6; i++ {
		fmt.Println(rl.Allow())
	}
}
