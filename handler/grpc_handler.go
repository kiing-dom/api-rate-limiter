package handler

import (
	"context"

	pb "github.com/kiing-dom/api-rate-limiter/proto"
	"github.com/kiing-dom/api-rate-limiter/store"
)

type GRPCRateLimiterServer struct {
	pb.UnimplementedRateLimiterServer
	store store.RLStore
}

func NewGRPCServer(s store.RLStore) *GRPCRateLimiterServer {
	return &GRPCRateLimiterServer{store: s}
}

func (g *GRPCRateLimiterServer) Check(
	ctx context.Context,
	req *pb.RateLimitRequest,
) (*pb.RateLimitResponse, error) {

	userID := req.GetUserId()
	if userID == "" {
		return &pb.RateLimitResponse{
			Allowed: false,
			Message: "Missing userID (X-API-KEY)",
		}, nil
	}

	algo := req.GetAlgo()
	rl := g.store.GetRateLimiter(userID, algo)

	if !rl.Allow(req.GetUserId()) {
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
