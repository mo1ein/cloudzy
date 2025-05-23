package weatherservice

import (
	"context"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"pricing/internal/domain"
	grpcadapter "pricing/internal/pkg"
	pb "pricing/proto/gen"
)

type Service struct {
	grpcClientAdapter *grpcadapter.GRPCClientAdapter
	logger            zerolog.Logger
}

func New(grpcClientAd *grpcadapter.GRPCClientAdapter, logger zerolog.Logger) *Service {
	return &Service{
		grpcClientAdapter: grpcClientAd,
		logger:            logger,
	}
}

func (s *Service) GetWeather(ctx context.Context) (*domain.Weather, error) {
	req := &pb.GetWeatherRequest{}
	rs, err := s.grpcClientAdapter.Do(func(conn *grpc.ClientConn) (any, error) {
		client := pb.NewWeatherServiceClient(conn)
		res, err := client.GetWeather(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, nil
	})

	if err != nil {
		s.logger.Err(err).Msg("failed to get weather data")
	}

	res, ok := rs.(*pb.GetWeatherResponse)
	if !ok {
		return nil, nil
	}

	weather := domain.Weather{
		Forecast:    res.Forecast,
		Temperature: res.Temperature,
		Altitude:    res.Altitude,
	}
	return &weather, nil
}
