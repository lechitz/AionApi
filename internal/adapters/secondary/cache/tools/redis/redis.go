// Package cache provides a Redis client for caching data.
package cache

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/redis/go-redis/v9"

	"github.com/lechitz/AionApi/internal/platform/config"
)

// FailedToConnectToRedis is a constant for logging errors when the Redis client fails to connect.
const FailedToConnectToRedis = "failed to connect to Redis"

// Client is an interface for caching operations, allowing setting, getting, and managing cached data.
type Client interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
}

// NewCacheConnection initializes a new Redis client using the provided configuration and logger.
// Returns the Redis client or an error if the connection fails.
func NewCacheConnection(cacheCfg config.CacheConfig, log output.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cacheCfg.Addr,
		Password: cacheCfg.Password,
		DB:       cacheCfg.DB,
		PoolSize: cacheCfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Errorw(FailedToConnectToRedis, "error", err)
		return nil, err
	}

	return client, nil
}
