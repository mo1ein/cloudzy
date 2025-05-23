package pricing

import (
	"context"
	"pricing/internal/constant"
)

func (s *Service) GetPrice(ctx context.Context) (float64, error) {
	price, err := s.repository.GetPricingHash(ctx, constant.RedisKey)
	if err != nil {
		return 0, err
	}
	return price, nil
}
