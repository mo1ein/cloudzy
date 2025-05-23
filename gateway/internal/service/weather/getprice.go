package weather

import (
	"context"
	"gateway/internal/domain"
	pb "github.com/mo1ein/cloudzy/proto/gen"
	"google.golang.org/grpc"
)

func (s *weatherService) GetWeather(ctx context.Context) (domain.Weather, error) {
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
		return domain.Weather{}, nil
	}

	return domain.Weather{
		Forecast:    res.Forecast,
		Temperature: res.Temperature,
		Altitude:    res.Altitude,
	}, nil
}
