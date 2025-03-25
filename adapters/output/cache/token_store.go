package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TokenStore struct {
	redis       *redis.Client
	loggerSugar *zap.SugaredLogger
}

func NewTokenStore(redisClient *redis.Client, logger *zap.SugaredLogger) cache.ITokenRepository {
	return &TokenStore{
		redis:       redisClient,
		loggerSugar: logger,
	}
}

func (ts *TokenStore) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := ts.formatTokenKey(token.UserID)

	expiration := 24 * time.Hour
	err := ts.redis.Set(ctx.BaseContext, key, token.Token, expiration).Err()
	if err != nil {
		ts.loggerSugar.Errorw("failed to save token to Redis", "key", key, "error", err)
		return err
	}

	ts.loggerSugar.Infow("token saved to Redis", "key", key)
	return nil
}

func (ts *TokenStore) Get(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	key := ts.formatTokenKey(token.UserID)

	value, err := ts.redis.Get(ctx.BaseContext, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("token not found for user ID %d", token.UserID)
		}
		ts.loggerSugar.Errorw("failed to get token from Redis", "key", key, "error", err)
		return "", err
	}

	return value, nil
}

func (ts *TokenStore) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := ts.formatTokenKey(token.UserID)

	expiration := 24 * time.Hour
	err := ts.redis.Set(ctx.BaseContext, key, token.Token, expiration).Err()
	if err != nil {
		ts.loggerSugar.Errorw("failed to update token in Redis", "key", key, "error", err)
		return err
	}

	ts.loggerSugar.Infow("token updated in Redis", "key", key)
	return nil
}

func (ts *TokenStore) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := ts.formatTokenKey(token.UserID)

	err := ts.redis.Del(ctx.BaseContext, key).Err()
	if err != nil {
		ts.loggerSugar.Errorw("failed to delete token from Redis", "key", key, "error", err)
		return err
	}

	ts.loggerSugar.Infow("token deleted from Redis", "key", key)
	return nil
}

func (ts *TokenStore) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
