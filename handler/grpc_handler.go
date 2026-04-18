package handler

import (
	"context"

	pb "github.com/kiing-dom/api-rate-limiter/proto"
	"github.com/kiing-dom/api-rate-limiter/store"
)

type GRPCRateLimiterServer struct {
	pb.UnimplementedRateLimiterServer
	store *store.Store
}

func NewGRPCServer(s *store.Store) *GRPCRateLimiterServer {
	return &GRPCRateLimiterServer{store: s}
}

func (g *GRPCRateLimiterServer) Check(
	ctx context.Context,
	req *pb.RateLimitRequest,
) (*pb.RateLimitResponse, error) {

	userID := req.GetUserId()

	rl := g.store.GetRateLimiter(userID)

	if !rl.Allow() {
		return &pb.RateLimitResponse{
			Allowed: false,
			Message: "Rate limited exceeded",
		}, nil
	}

	return &pb.RateLimitResponse{
		Allowed: true,
		Message: "Request allowed",
	}, nil
}
