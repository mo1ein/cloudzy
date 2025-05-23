package getservice

import (
	"context"
	"weather/internal/constant"
	"weather/internal/domain"
)

func (s *Service) GetWeather(ctx context.Context) (domain.Weather, error) {
	weather, err := s.repo.GetWeatherHash(ctx, constant.RedisKey)
	if err != nil {
		return domain.Weather{}, err
	}
	return *weather, nil
}
