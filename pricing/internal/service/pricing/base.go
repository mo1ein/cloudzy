package pricing

import (
	"context"
	"pricing/internal/domain"
)

type weatherService interface {
	GetWeather(ctx context.Context) (*domain.Weather, error)
}

type repository interface {
	GetPricingHash(ctx context.Context, key string) (float64, error)
}

type Service struct {
	weatherService weatherService
	repository     repository
}

func New(weatherService weatherService, repository repository) *Service {
	return &Service{
		weatherService: weatherService,
		repository:     repository,
	}
}
