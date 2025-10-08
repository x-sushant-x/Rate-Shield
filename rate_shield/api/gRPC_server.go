package api

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/limiter"
	"github.com/x-sushant-x/RateShield/proto/github.com/x-sushant-x/RateShield/ratelimitpb"
	"github.com/x-sushant-x/RateShield/utils"
	"google.golang.org/grpc"
)

type gRPCService struct {
	ratelimitpb.UnimplementedRateLimitServiceServer
	limiterSvc *limiter.Limiter
}

func newgRPCService(limiterSvc *limiter.Limiter) *gRPCService {
	return &gRPCService{
		limiterSvc: limiterSvc,
	}
}

func (s *gRPCService) CheckRateLimit(ctx context.Context, req *ratelimitpb.RateLimitRequest) (*ratelimitpb.RateLimitResponse, error) {
	ip := req.GetIp()
	endpoint := req.GetEndpoint()

	if err := utils.ValidateLimitRequest(req.Ip, req.Endpoint); err != nil {
		return &ratelimitpb.RateLimitResponse{
			HttpStatusCode: 400,
		}, nil
	}

	resp := s.limiterSvc.CheckLimit(ip, endpoint)

	return &ratelimitpb.RateLimitResponse{
		HttpStatusCode: int32(resp.HTTPStatusCode),
		Limit:          int32(resp.RateLimit_Limit),
		Remaining:      int32(resp.RateLimit_Remaining),
	}, nil
}

func StartGRPCServer(limiterSvc *limiter.Limiter, port string) {
	grpcServer := grpc.NewServer()

	grpcService := newgRPCService(limiterSvc)
	ratelimitpb.RegisterRateLimitServiceServer(grpcServer, grpcService)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg("gRPC server listening on :" + port + " ✅")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal().Err(err)
	}
}
