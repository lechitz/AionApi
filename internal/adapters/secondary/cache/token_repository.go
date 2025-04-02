package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	cacheport "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TokenRepository struct {
	cache  cacheport.Store
	logger *zap.SugaredLogger
}

func NewTokenRepository(cache cacheport.Store, logger *zap.SugaredLogger) *TokenRepository {
	return &TokenRepository{
		cache:  cache,
		logger: logger,
	}
}

func (ts *TokenRepository) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := ts.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour
	err := ts.cache.Set(key, token.Token, expiration)
	if err != nil {
		ts.logger.Errorw(constants.ErrorToSaveTokenToRedis, constants.Key, key, constants.Error, err)
		return err
	}

	ts.logger.Infow(constants.SuccessToSaveTokenToRedis, constants.Key, key)
	return nil
}

func (ts *TokenRepository) Get(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	key := ts.formatTokenKey(token.UserID)

	value, err := ts.cache.Get(key)
	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "redis: nil" {
			return "", fmt.Errorf("token not found for user ID %d", token.UserID)
		}
		ts.logger.Errorw(constants.ErrorToGetTokenFromRedis, constants.Key, key, constants.Error, err)
		return "", err
	}

	ts.logger.Infow(constants.SuccessToGetTokenFromRedis, constants.Key, key)
	return value, nil
}

func (ts *TokenRepository) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := ts.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour
	err := ts.cache.Set(key, token.Token, expiration)
	if err != nil {
		ts.logger.Errorw(constants.ErrorToUpdateTokenInRedis, constants.Key, key, constants.Error, err)
		return err
	}

	ts.logger.Infow(constants.SuccessToUpdateTokenInRedis, constants.Key, key)
	return nil
}

func (ts *TokenRepository) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	key := ts.formatTokenKey(token.UserID)

	err := ts.cache.Delete(key)
	if err != nil {
		ts.logger.Errorw(constants.ErrorToDeleteTokenFromRedis, constants.Key, key, constants.Error, err)
		return err
	}

	ts.logger.Infow(constants.SuccessToDeleteTokenFromRedis, constants.Key, key)
	return nil
}

func (ts *TokenRepository) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
