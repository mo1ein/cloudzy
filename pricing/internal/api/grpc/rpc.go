package grpcapi

import (
	"context"
	pb "pricing/proto/gen"
)

func (s *Server) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	return s.handler.GetPrice(ctx, req)
}
