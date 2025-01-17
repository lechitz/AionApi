package redis

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/core/domain/entities"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type CacheRepo struct {
	redis       *redis.Client
	loggerSugar *zap.SugaredLogger
}

var _ cache.ITokenRepository = (*CacheRepo)(nil)

func NewCacheRepo(redisClient *redis.Client, loggerSugar *zap.SugaredLogger) cache.ITokenRepository {
	return &CacheRepo{
		redis:       redisClient,
		loggerSugar: loggerSugar,
	}
}

func (c *CacheRepo) SaveToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error {
	redisCtx := ctx.BaseContext

	//TODO: Move expiration to config
	expiration := 24 * time.Hour

	key := c.formatTokenKey(tokenDomain.UserID)
	return c.redis.Set(redisCtx, key, tokenDomain.Token, expiration).Err()
}

func (c *CacheRepo) GetToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) (string, error) {
	redisCtx := ctx.BaseContext
	key := c.formatTokenKey(tokenDomain.UserID)

	token, err := c.redis.Get(redisCtx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf(ErrorTokenNotFound)
		}
		return "", err
	}
	return token, nil
}

func (c *CacheRepo) UpdateToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error {
	return c.SaveToken(ctx, tokenDomain)
}

func (c *CacheRepo) DeleteToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error {
	redisCtx := ctx.BaseContext
	key := c.formatTokenKey(tokenDomain.UserID)
	return c.redis.Del(redisCtx, key).Err()
}

// Auxiliar functions

func (c *CacheRepo) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
