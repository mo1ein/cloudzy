package grpctransformes

import (
	pb "github.com/mo1ein/cloudzy/proto/gen"
	"weather/internal/domain"
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
