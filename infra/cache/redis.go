package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RedisClient struct {
	Client      *redis.Client
	LoggerSugar *zap.SugaredLogger
}

func NewRedisClient(addr, password string, db int, loggerSugar *zap.SugaredLogger) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		loggerSugar.Fatalf("Failed to connect to Redis: %v", err)
	}

	loggerSugar.Infow("redis connection established", "address", addr, "db", db)

	return &RedisClient{
		Client:      client,
		LoggerSugar: loggerSugar,
	}
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
