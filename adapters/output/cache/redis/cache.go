package redis

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type CacheRepo struct {
	redis       *redis.Client
	loggerSugar *zap.SugaredLogger
}

func NewCacheRepo(redisClient *redis.Client, loggerSugar *zap.SugaredLogger) cache.ITokenRepository {
	return &CacheRepo{
		redis:       redisClient,
		loggerSugar: loggerSugar,
	}
}

func (c *CacheRepo) SaveToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error {
	redisCtx := ctx.BaseContext

	key := c.formatTokenKey(tokenDomain.UserID)

	c.loggerSugar.Infow("Saving token to Redis", "key", key, "token", tokenDomain.Token)

	expiration := 24 * time.Hour
	err := c.redis.Set(redisCtx, key, tokenDomain.Token, expiration).Err()
	if err != nil {
		c.loggerSugar.Errorw("Failed to save token", "key", key, "error", err)
		return err
	}
	return nil
}

func (c *CacheRepo) GetToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) (string, error) {
	redisCtx := ctx.BaseContext
	key := c.formatTokenKey(tokenDomain.UserID)

	token, err := c.redis.Get(redisCtx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("token not found for user ID %d", tokenDomain.UserID)
		}
		return "", err
	}
	return token, nil
}

func (c *CacheRepo) UpdateToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error {

	redisCtx := ctx.BaseContext
	key := c.formatTokenKey(tokenDomain.UserID)
	expiration := 24 * time.Hour

	return c.redis.Set(redisCtx, key, tokenDomain.Token, expiration).Err()
}

func (c *CacheRepo) DeleteToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error {
	redisCtx := ctx.BaseContext
	key := c.formatTokenKey(tokenDomain.UserID)
	return c.redis.Del(redisCtx, key).Err()
}

func (c *CacheRepo) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
