package main

import (
	"log"
	"net/http"

	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
)

func main() {
	http.HandleFunc("/", rate_limiter.RateLimitHandler)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Server failed to run %v", err)
	}
}
