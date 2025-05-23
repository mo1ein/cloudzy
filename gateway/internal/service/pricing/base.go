package pricing

import (
	grpcadapter "gateway/internal/pkg"
	"github.com/rs/zerolog"
)

type pricingService struct {
	logger            zerolog.Logger
	grpcClientAdapter *grpcadapter.GRPCClientAdapter
}

func New(grpcClientAdapter *grpcadapter.GRPCClientAdapter, logger zerolog.Logger) *pricingService {
	return &pricingService{
		grpcClientAdapter: grpcClientAdapter,
		logger:            logger,
	}
}
