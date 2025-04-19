package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config"
)

const FailedToConnectToRedis = "failed to connect to Redis"

type CacheClient interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
}

func NewCacheConnection(cfg config.CacheConfig, log logger.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Errorw(FailedToConnectToRedis, "error", err)
		return nil, err
	}

	return client, nil
}
