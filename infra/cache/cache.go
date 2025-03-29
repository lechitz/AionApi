package cache

import (
	"context"
	"github.com/lechitz/AionApi/internal/platform/config"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	FailedToConnectToRedis  = "Failed to connect to Redis: %v"
	SuccessToConnectToRedis = "Redis connection established"
)

func NewCacheConnection(config config.CacheConfig, loggerSugar *zap.SugaredLogger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		loggerSugar.Fatalf(FailedToConnectToRedis, err)
	}

	loggerSugar.Infow(SuccessToConnectToRedis, "address", config.Addr, "db", config.DB)

	return client
}
