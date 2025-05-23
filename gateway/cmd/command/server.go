package command

import (
	"context"
	"fmt"
	"gateway/internal/api/rest"
	get "gateway/internal/api/rest/handler"
	"gateway/internal/config"
	grpcadapter "gateway/internal/pkg"
	"gateway/internal/service/pricing"
	"gateway/internal/service/weather"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type Server struct {
	Logger zerolog.Logger
}

func (cmd *Server) Command(ctx context.Context, logger zerolog.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "run gateway server",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.Logger = logger
			cmd.main(cfg, ctx)
		},
	}
}

func (cmd *Server) main(cfg *config.Config, ctx context.Context) {

	pricingGrpcClientAdapter, err := grpcadapter.New(cfg.PricingSvc.BaseURL, &cmd.Logger)
	if err != nil {
		cmd.Logger.Error().
			Err(err).
			Msg("failed to initialize grpc pricing adapter")
		return
	}

	weatherGrpcClientAdapter, err := grpcadapter.New(cfg.WeatherSvc.BaseURL, &cmd.Logger)
	if err != nil {
		cmd.Logger.Error().
			Err(err).
			Msg("failed to initialize grpc pricing adapter")
		return
	}

	pricingSvc := pricing.New(pricingGrpcClientAdapter, cmd.Logger)
	weatherSvc := weather.New(weatherGrpcClientAdapter, cmd.Logger)
	getHandler := get.NewHandler(pricingSvc, weatherSvc)

	server := rest.New()
	server.SetupAPIRoutes(getHandler)

	if err := server.Serve(ctx, fmt.Sprintf("%s:%d", cfg.HTTP.APIHost, cfg.HTTP.APIPort)); err != nil {
		cmd.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
