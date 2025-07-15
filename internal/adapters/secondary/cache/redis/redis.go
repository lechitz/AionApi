// Package redis provides a Redis cache implementation.
package redis

import (
	"context"
	"errors"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/redis/go-redis/v9"
)

// FailedToConnectToRedis is the error message for when the Redis client fails to connect.
const FailedToConnectToRedis = "failed to connect to Redis"

type redisClient struct {
	client *redis.Client
	logger output.ContextLogger
}

// NewConnection initializes a new Redis cache connection.
func NewConnection(appCtx context.Context, cfg config.CacheConfig, logger output.ContextLogger) (output.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(appCtx, cfg.ConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Errorw(FailedToConnectToRedis, commonkeys.Error, err)
		return nil, err
	}

	return &redisClient{client: client, logger: logger}, nil
}

func (r *redisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (r *redisClient) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *redisClient) Close() error {
	return r.client.Close()
}
