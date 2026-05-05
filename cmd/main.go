package main

import (
	"log"
	"net/http"

	"github.com/kiing-dom/api-rate-limiter/handler"
	"github.com/kiing-dom/api-rate-limiter/internal/config"
	"github.com/kiing-dom/api-rate-limiter/internal/server"
	"github.com/kiing-dom/api-rate-limiter/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	store, err := store.NewStore(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	// gRPC
	go server.StartGRPCServer(store)

	// http
	h := handler.RateLimitHandler(store)
	http.HandleFunc("/", h)

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Server failed to run %v", err)
	}
}
