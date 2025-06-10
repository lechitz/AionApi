package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/adapters/secondary/cache/constants"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	cache  *redis.Client
	logger logger.Logger
}

func NewTokenRepository(cache *redis.Client, logger logger.Logger) *TokenRepository {
	return &TokenRepository{
		cache:  cache,
		logger: logger,
	}
}

func (t *TokenRepository) Save(ctx context.Context, token domain.TokenDomain) error {
	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour

	if err := t.cache.Set(ctx, key, token.Token, expiration).Err(); err != nil {
		t.logger.Errorw(constants.ErrorToSaveTokenToRedis, constants.Key, key, constants.Error, err)
		return err
	}

	return nil
}

func (t *TokenRepository) Get(ctx context.Context, token domain.TokenDomain) (string, error) {
	key := t.formatTokenKey(token.UserID)

	value, err := t.cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "redis: nil" {
			return "", fmt.Errorf("token not found for user ID %d", token.UserID)
		}
		t.logger.Errorw(
			constants.ErrorToGetTokenFromRedis,
			constants.Key,
			key,
			constants.Error,
			err,
		)
		return "", err
	}

	return value, nil
}

func (t *TokenRepository) Update(ctx context.Context, token domain.TokenDomain) error {
	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour

	if err := t.cache.Set(ctx, key, token.Token, expiration).Err(); err != nil {
		t.logger.Errorw(
			constants.ErrorToUpdateTokenInRedis,
			constants.Key,
			key,
			constants.Error,
			err,
		)
		return err
	}

	t.logger.Infow(constants.SuccessToUpdateTokenInRedis, constants.Key, key)
	return nil
}

func (t *TokenRepository) Delete(ctx context.Context, token domain.TokenDomain) error {
	key := t.formatTokenKey(token.UserID)

	if err := t.cache.Del(ctx, key).Err(); err != nil {
		t.logger.Errorw(
			constants.ErrorToDeleteTokenFromRedis,
			constants.Key,
			key,
			constants.Error,
			err,
		)
		return err
	}

	return nil
}

func (t *TokenRepository) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
