package command

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"pricing/internal/config"
	"pricing/internal/constant"
	grpcadapter "pricing/internal/pkg"
	"pricing/internal/repository"
	"pricing/internal/service/pricing"
	weatherSvc "pricing/internal/service/weather"
	"time"
)

type Worker struct {
	logger zerolog.Logger
}

func (cmd Worker) Command(ctx context.Context, logger zerolog.Logger, cfg *config.Config) *cobra.Command {

	return &cobra.Command{
		Use:   "pricing-worker",
		Short: "run pricing worker",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.logger = logger
			cmd.main(ctx, cfg)
		},
	}
}

func (cmd *Worker) main(ctx context.Context, cfg *config.Config) {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.Database,
	})
	defer redisConn.Close()

	repo := repository.New(redisConn)
	weatherGrpcClientAdapter, err := grpcadapter.New(cfg.WeatherSvc.BaseURL, &cmd.logger)
	if err != nil {
		cmd.logger.Error().
			Err(err).
			Msg("failed to initialize grpc adapter")
		return
	}
	weatherSvc := weatherSvc.New(weatherGrpcClientAdapter, cmd.logger)
	pricingSvc := pricing.New(weatherSvc, repo)

	ticker := time.NewTicker(1 * time.Second)
	cmd.logger.Info().
		Msg("pricing:starting worker...")

	for {
		select {
		case <-ticker.C:
			price, err := pricingSvc.Calculate(ctx)
			if err != nil {
				cmd.logger.Err(err).Str("method", "worker").Msg("failed to calculate price")
				continue
			}
			// todo: ttl?
			err = repo.UpdateWeatherHash(ctx, constant.RedisKey, price, 1*time.Hour)
			if err != nil {
				cmd.logger.Error().Err(err).Str("method", "worker").Msg("failed to set weather data")
				continue
			}
			cmd.logger.Info().
				Str("key", constant.RedisKey).
				Float64("price", price). // Assuming Temperature is float64
				Msg("updated price hash")

		case <-ctx.Done():
			return
		}
	}
}
