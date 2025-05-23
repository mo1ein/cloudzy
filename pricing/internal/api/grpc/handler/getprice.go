package grpchandlers

import (
	"context"
	pb "pricing/proto/gen"
)

type pricingService interface {
	GetPrice(ctx context.Context) (float64, error)
}

type Handler struct {
	pricingSvc pricingService
}

func New(
	pricingSvc pricingService,
) Handler {
	return Handler{
		pricingSvc: pricingSvc,
	}
}

func (h Handler) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	price, err := h.pricingSvc.GetPrice(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetPriceResponse{
		Price: price,
	}, nil
}
