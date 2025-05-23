package grpctransformes

import (
	"weather/internal/domain"
	pb "weather/proto/gen"
)

type Transformer struct {
}

func New() Transformer {
	return Transformer{}
}

func (r *Transformer) Transform(weather domain.Weather) *pb.GetWeatherResponse {
	return &pb.GetWeatherResponse{
		Forecast:    weather.Forecast,
		Temperature: weather.Temperature,
		Altitude:    weather.Altitude,
	}
}
