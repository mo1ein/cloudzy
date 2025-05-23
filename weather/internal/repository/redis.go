package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"weather/internal/domain"
)

type Repository struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *Repository {
	return &Repository{
		redisClient: redisClient,
	}
}

func (r *Repository) UpdateWeatherHash(ctx context.Context, key string, data domain.Weather, ttl time.Duration) error {
	// Use pipeline for atomic operations
	pipe := r.redisClient.TxPipeline()

	pipe.HSet(ctx, key,
		"temp", data.Temperature,
		"altitude", data.Altitude,
		"forecast", data.Forecast,
	)

	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) GetWeatherHash(ctx context.Context, key string) (*domain.Weather, error) {
	result, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash: %w", err)
	}

	// Convert temperature string to float64
	temp, err := strconv.ParseFloat(result["temp"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid temperature format: %w", err)
	}

	// Convert altitude string to float64
	altitude, err := strconv.ParseFloat(result["altitude"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid altitude format: %w", err)
	}

	return &domain.Weather{
		Temperature: temp,
		Altitude:    altitude,
		Forecast:    result["forecast"],
	}, nil
}

func (r *Repository) CreateHSet(ctx context.Context, key string, values ...interface{}) error {
	err := r.redisClient.HSet(ctx, key, values).Err()
	if err != nil {
		return err
	}
	// todo: expire things

	return nil
}
