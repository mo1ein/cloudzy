package get

import (
	"context"
	"gateway/internal/domain"
)

type pricingSvc interface {
	GetPrice(ctx context.Context) (float64, error)
}

type weatherSvc interface {
	GetWeather(ctx context.Context) (domain.Weather, error)
}

type Handler struct {
	pricingSvc pricingSvc
	weatherSvc weatherSvc
}

func NewHandler(pricingSvc pricingSvc, weatherSvc weatherSvc) Handler {
	return Handler{
		pricingSvc: pricingSvc,
		weatherSvc: weatherSvc,
	}
}
