package command

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"net"
	grpcadapter "pricing/internal/pkg"
	"pricing/internal/repository"
	"pricing/internal/service/pricing"
	weatherservice "pricing/internal/service/weather"

	grpcapi "pricing/internal/api/grpc"
	grpchandlers "pricing/internal/api/grpc/handler"
	"pricing/internal/config"
)

type GRPCServer struct {
	Logger zerolog.Logger
}

func (cmd *GRPCServer) Command(ctx context.Context, logger zerolog.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc-server",
		Short: "run post grpc server",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.Logger = logger
			cmd.main(ctx, cfg)
		},
	}
}

func (cmd *GRPCServer) main(ctx context.Context, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		cmd.Logger.Error().
			Str("addr", addr).
			Err(err).
			Msg("failed to listen")
		return
	}

	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.Database,
	})
	defer redisConn.Close()

	weatherGrpcClientAdapter, err := grpcadapter.New(cfg.WeatherSvc.BaseURL, &cmd.Logger)
	if err != nil {
		cmd.Logger.Error().
			Err(err).
			Msg("failed to initialize grpc adapter")
		return
	}

	weatherService := weatherservice.New(weatherGrpcClientAdapter, cmd.Logger)
	repo := repository.New(redisConn)
	pricingSvc := pricing.New(weatherService, repo)
	grpcHandler := grpchandlers.New(pricingSvc)

	// Wire up the gRPC server, passing in your zerolog.Logger
	grpcServer := grpcapi.New(listener, grpcHandler, cmd.Logger)

	cmd.Logger.Info().
		Str("addr", addr).
		Msg("starting grpc server")

	if err := grpcServer.Start(ctx); err != nil {
		cmd.Logger.Error().
			Err(err).
			Msg("failed to start grpc server")
	}
}
