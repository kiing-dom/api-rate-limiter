package rate_limiter

import (
	"log"
	"net"
	"net/http"
)

func RateLimitHandler(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error parsing host port %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Internal Server Error"))
	}

	userID := host
	rl := GetRateLimiter(userID, 5, 1)

	if !rl.Allow() {
		log.Printf("Too Many Request by user: %s. Try again later", userID)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Rate limit exceeded"))
		return
	}

	w.Write([]byte("Request allowed!"))
}
