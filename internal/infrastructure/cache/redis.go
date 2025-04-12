package cache

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/redis/go-redis/v9"
)

type CacheClient interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

func NewCacheConnection(cfg config.CacheConfig, logger logger.Logger) *redis.Client {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Errorf("Failed to connect to Redis: %v", err)
	}

	logger.Infow("Redis connected", "addr", cfg.Addr, "db", cfg.DB)

	return client
}
