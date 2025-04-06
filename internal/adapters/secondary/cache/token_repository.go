package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TokenRepository struct {
	cache  *redis.Client
	logger *zap.SugaredLogger
}

func NewTokenRepository(cache *redis.Client, loggerSugar *zap.SugaredLogger) *TokenRepository {
	return &TokenRepository{
		cache:  cache,
		logger: loggerSugar,
	}
}

//func (t *TokenRepository) CreateToken(ctx domain.ContextControl, token domain.TokenDomain) error {
//	key := t.formatTokenKey(token.UserID)
//	expiration := 24 * time.Hour
//	err := t.cache.Set(ctx.BaseContext, key, token.Token, expiration).Err()
//	if err != nil {
//		t.logger.Errorw(constants.ErrorToCreateTokenInRedis, constants.Key, key, constants.Error, err)
//		return err
//	}
//
//	t.logger.Infow(constants.SuccessToCreateTokenInRedis, constants.Key, key)
//	return nil
//}

func (t *TokenRepository) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour
	err := t.cache.Set(ctx.BaseContext, key, token.Token, expiration).Err()
	if err != nil {
		t.logger.Errorw(constants.ErrorToSaveTokenToRedis, constants.Key, key, constants.Error, err)
		return err
	}

	t.logger.Infow(constants.SuccessToSaveTokenToRedis, constants.Key, key)
	return nil
}

func (t *TokenRepository) Get(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	key := t.formatTokenKey(token.UserID)

	value, err := t.cache.Get(ctx.BaseContext, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "redis: nil" {
			return "", fmt.Errorf("token not found for user ID %d", token.UserID)
		}
		t.logger.Errorw(constants.ErrorToGetTokenFromRedis, constants.Key, key, constants.Error, err)
		return "", err
	}

	t.logger.Infow(constants.SuccessToGetTokenFromRedis, constants.Key, key)
	return value, nil
}

func (t *TokenRepository) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour
	err := t.cache.Set(ctx.BaseContext, key, token.Token, expiration).Err()
	if err != nil {
		t.logger.Errorw(constants.ErrorToUpdateTokenInRedis, constants.Key, key, constants.Error, err)
		return err
	}

	t.logger.Infow(constants.SuccessToUpdateTokenInRedis, constants.Key, key)
	return nil
}

func (t *TokenRepository) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := t.formatTokenKey(token.UserID)

	err := t.cache.Del(ctx.BaseContext, key).Err()
	if err != nil {
		t.logger.Errorw(constants.ErrorToDeleteTokenFromRedis, constants.Key, key, constants.Error, err)
		return err
	}

	return nil
}

func (t *TokenRepository) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
