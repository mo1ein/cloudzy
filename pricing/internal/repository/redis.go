package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type Repository struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *Repository {
	return &Repository{
		redisClient: redisClient,
	}
}

func (r *Repository) UpdateWeatherHash(ctx context.Context, key string, data float64, ttl time.Duration) error {
	// Use pipeline for atomic operations
	pipe := r.redisClient.TxPipeline()

	pipe.HSet(ctx, key,
		"price", data,
	)

	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) GetPricingHash(ctx context.Context, key string) (float64, error) {
	result, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get hash: %w", err)
	}

	// Convert temperature string to float64
	price, err := strconv.ParseFloat(result["price"], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid temperature format: %w", err)
	}

	return price, nil
}

func (r *Repository) CreateHSet(ctx context.Context, key string, values ...interface{}) error {
	err := r.redisClient.HSet(ctx, key, values).Err()
	if err != nil {
		return err
	}
	// todo: expire things

	return nil
}
