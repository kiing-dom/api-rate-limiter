package main

import (
	"log"
	"net/http"

	"github.com/kiing-dom/api-rate-limiter/handler"
	"github.com/kiing-dom/api-rate-limiter/store"
)

func main() {
	store := store.NewStore()
	h := handler.RateLimitHandler(store)
	http.HandleFunc("/", h)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Server failed to run %v", err)
	}
}
