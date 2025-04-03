package cache

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/config"
	cacheport "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCache struct {
	client *redis.Client
	logger *zap.SugaredLogger
	ctx    context.Context
}

func NewRedisConnection(cfg config.CacheConfig, logger *zap.SugaredLogger) cacheport.Store {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger.Infow("Redis connected", "addr", cfg.Addr, "db", cfg.DB)

	return &RedisCache{
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}

func (r *RedisCache) Set(key string, value any, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
