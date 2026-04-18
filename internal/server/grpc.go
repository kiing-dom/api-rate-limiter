package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/kiing-dom/api-rate-limiter/handler"
	pb "github.com/kiing-dom/api-rate-limiter/proto"
	"github.com/kiing-dom/api-rate-limiter/store"
)

func StartGRPCServer(s *store.Store) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRateLimiterServer(
		grpcServer,
		handler.NewGRPCServer(s),
	)

	log.Println("gRPC Server running on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
