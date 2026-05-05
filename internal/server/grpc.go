package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/kiing-dom/api-rate-limiter/handler"
	"github.com/kiing-dom/api-rate-limiter/internal/config"
	pb "github.com/kiing-dom/api-rate-limiter/proto"
	"github.com/kiing-dom/api-rate-limiter/store"
)

func StartGRPCServer(s *store.Store, cfg *config.Config) {
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRateLimiterServer(
		grpcServer,
		handler.NewGRPCServer(s),
	)

	log.Printf("gRPC Server running on +:%s", cfg.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
