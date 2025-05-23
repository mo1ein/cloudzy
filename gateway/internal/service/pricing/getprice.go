package pricing

import (
	"context"
	//_ pb "gateway/proto/gen"
	pb "github.com/mo1ein/cloudzy/proto/gen"
	"google.golang.org/grpc"
)

func (s *pricingService) GetPrice(ctx context.Context) (float64, error) {
	req := &pb.GetPriceRequest{}
	rs, err := s.grpcClientAdapter.Do(func(conn *grpc.ClientConn) (any, error) {
		client := pb.NewWeatherServiceClient(conn)
		res, err := client.GetPrice(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, nil
	})

	if err != nil {
		s.logger.Err(err).Msg("failed to get weather data")
	}

	res, ok := rs.(*pb.GetPriceResponse)
	if !ok {
		return 0, nil
	}

	return res.Price, nil
}
