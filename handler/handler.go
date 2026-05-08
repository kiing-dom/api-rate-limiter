package handler

import (
	"log"
	"net/http"

	"github.com/kiing-dom/api-rate-limiter/store"
)

func RateLimitHandler(s store.RLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-API-Key")
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing X-API-Key header"))
			return
		}

		algo := r.URL.Query().Get("algo")
		rl := s.GetRateLimiter(userID, algo)

		if !rl.Allow(userID) {
			log.Printf("Too Many Request by user: %s. Try again later", userID)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded"))
			return
		}

		w.Write([]byte("Request allowed!"))
	}
}
