package command

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"net/http"
	"time"
	"weather/internal/config"
	"weather/internal/constant"
	"weather/internal/repository"
	getservice "weather/internal/service"
)

type Worker struct {
	logger zerolog.Logger
}

func (cmd Worker) Command(ctx context.Context, logger zerolog.Logger, cfg *config.Config) *cobra.Command {

	return &cobra.Command{
		Use:   "weather-worker",
		Short: "run weather worker",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.logger = logger
			cmd.main(ctx, cfg)
		},
	}
}

// todo: fix this env
// todo: add random lat & lon
const lat = 68
const lon = 34

func (cmd *Worker) main(ctx context.Context, cfg *config.Config) {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.Database,
	})
	defer redisConn.Close()

	repo := repository.New(redisConn)
	weatherSvc := getservice.New(&http.Client{}, *repo)

	ticker := time.NewTicker(1 * time.Second)
	cmd.logger.Info().
		Msg("starting worker...")

	for {
		select {
		case <-ticker.C:
			data, err := weatherSvc.FetchWeather(ctx, lat, lon)
			if err != nil {
				cmd.logger.Error().
					Err(err).
					Str("method", "worker").
					Msg("failed to fetch weather data")
				continue
			}
			// todo: add time to .env
			// todo: ttl?
			err = repo.UpdateWeatherHash(ctx, constant.RedisKey, data, 1*time.Hour)
			if err != nil {
				cmd.logger.Error().Err(err).Str("method", "worker").Msg("failed to set weather data")
				continue
			}
			cmd.logger.Info().
				Str("key", constant.RedisKey).
				Float64("temp", data.Temperature).  // Assuming Temperature is float64
				Float64("altitude", data.Altitude). // Assuming Altitude is int
				Str("forecast", data.Forecast).
				Msg("updated weather hash")

		case <-ctx.Done():
			return
		}
	}
}
