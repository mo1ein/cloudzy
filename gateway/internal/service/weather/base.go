package weather

import (
	grpcadapter "gateway/internal/pkg"
	"github.com/rs/zerolog"
)

type weatherService struct {
	logger            zerolog.Logger
	grpcClientAdapter *grpcadapter.GRPCClientAdapter
}

func New(grpcClientAdapter *grpcadapter.GRPCClientAdapter, logger zerolog.Logger) *weatherService {
	return &weatherService{
		grpcClientAdapter: grpcClientAdapter,
		logger:            logger,
	}
}
