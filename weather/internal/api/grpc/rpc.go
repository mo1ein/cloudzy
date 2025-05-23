package grpcapi

import (
	"context"
	pb "github.com/mo1ein/cloudzy/proto/gen"
)

func (s *Server) GetWeather(ctx context.Context, req *pb.GetWeatherRequest) (*pb.GetWeatherResponse, error) {
	return s.handler.GetWeather(ctx)
}
