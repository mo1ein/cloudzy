package grpcapi

import (
	"context"
	pb "weather/proto/gen"
)

func (s *Server) GetWeather(ctx context.Context, req *pb.GetWeatherRequest) (*pb.GetWeatherResponse, error) {
	return s.handler.GetWeather(ctx)
}
