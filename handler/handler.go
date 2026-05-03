package handler

import (
	"log"
	"net"
	"net/http"

	"github.com/kiing-dom/api-rate-limiter/store"
)

func RateLimitHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error parsing host port %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Internal Server Error"))
			return
		}

		userID := host
		rl := s.GetRateLimiter(userID)

		if !rl.Allow(userID) {
			log.Printf("Too Many Request by user: %s. Try again later", userID)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded"))
			return
		}

		w.Write([]byte("Request allowed!"))
	}
}
