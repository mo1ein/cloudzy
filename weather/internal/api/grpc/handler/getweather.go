package grpchandlers

import (
	"context"
	grpctransformes "weather/internal/api/grpc/transformer"
	"weather/internal/domain"
	pb "weather/proto/gen"
)

type weatherService interface {
	FetchWeather(ctx context.Context, lat, lon float64) (domain.Weather, error)
	GetWeather(ctx context.Context) (domain.Weather, error)
}

type Handler struct {
	weatherSvc   weatherService
	weatherTrans grpctransformes.Transformer
}

func New(
	weatherService weatherService,
	weatherTrans grpctransformes.Transformer,
) Handler {
	return Handler{
		weatherSvc:   weatherService,
		weatherTrans: weatherTrans,
	}
}

func (h Handler) GetWeather(ctx context.Context) (*pb.GetWeatherResponse, error) {
	weather, err := h.weatherSvc.GetWeather(ctx)
	if err != nil {
		return nil, err
	}
	return h.weatherTrans.Transform(weather), nil
}
