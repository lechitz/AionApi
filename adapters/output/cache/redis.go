package cache

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/infra/cache"
	"github.com/lechitz/AionApi/internal/core/constants"
	"go.uber.org/zap"
	"time"
)

type TokenStore struct {
	RedisClient *cache.RedisClient
	LoggerSugar *zap.SugaredLogger
	SecretKey   []byte
}

func NewTokenStore(redisClient *cache.RedisClient, loggerSugar *zap.SugaredLogger, secretKey []byte) *TokenStore {
	return &TokenStore{
		RedisClient: redisClient,
		LoggerSugar: loggerSugar,
		SecretKey:   secretKey,
	}
}

func (t *TokenStore) CreateToken(ctx context.Context, userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(t.SecretKey)
	if err != nil {
		t.LoggerSugar.Errorw(constants.ErrorToCreateToken, "error", err.Error())
		return "", err
	}

	key := fmt.Sprintf("user:%d:token", userID)
	err = t.RedisClient.Client.Set(ctx, key, signedToken, 24*time.Hour).Err()
	if err != nil {
		t.LoggerSugar.Errorw(constants.ErrorToStoreTokenInRedis, "error", err.Error(), "userID", userID)
		return "", err
	}

	t.LoggerSugar.Infow(constants.SuccessTokenCreated, "userID", userID)
	return signedToken, nil
}

func (t *TokenStore) ValidateToken(ctx context.Context, tokenFromCookie string) (string, uint64, error) {
	parsedToken, err := jwt.Parse(tokenFromCookie, func(token *jwt.Token) (interface{}, error) {
		return t.SecretKey, nil
	})
	if err != nil || !parsedToken.Valid {
		t.LoggerSugar.Errorw(constants.ErrorInvalidToken, "token", tokenFromCookie, "error", err)
		return "", 0, fmt.Errorf(constants.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.LoggerSugar.Errorw(constants.ErrorInvalidTokenClaims, "token", tokenFromCookie)
		return "", 0, fmt.Errorf(constants.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		t.LoggerSugar.Errorw(constants.ErrorInvalidUserIDClaim, "token", tokenFromCookie)
		return "", 0, fmt.Errorf(constants.ErrorInvalidUserIDClaim)
	}

	userID := uint64(userIDFloat)

	tokenFromRedis, err := t.GetTokenByUserID(ctx, userID)
	if err != nil {
		t.LoggerSugar.Errorw(constants.ErrorToRetrieveTokenFromRedis, "error", err.Error(), "userID", userID)
		return "", 0, fmt.Errorf(constants.ErrorToRetrieveTokenFromRedis)
	}

	if tokenFromCookie != tokenFromRedis {
		t.LoggerSugar.Errorw(constants.ErrorTokenMismatch, "userID", userID, "tokenFromCookie", tokenFromCookie, "tokenFromRedis", tokenFromRedis)
		return "", 0, fmt.Errorf(constants.ErrorTokenMismatch)
	}

	t.LoggerSugar.Infow(constants.SuccessTokenValidated, "userID", userID)
	return tokenFromCookie, userID, nil
}

func (t *TokenStore) GetTokenByUserID(ctx context.Context, userID uint64) (string, error) {
	key := fmt.Sprintf("user:%d:token", userID)
	value, err := t.RedisClient.Client.Get(ctx, key).Result()
	if err != nil {
		t.LoggerSugar.Errorw(constants.ErrorToRetrieveTokenFromRedis, "error", err.Error(), "key", key)
		return "", err
	}

	t.LoggerSugar.Infow(constants.InfoTokenRetrieved, "userID", userID)
	return value, nil
}

func (t *TokenStore) DeleteTokenByUserID(ctx context.Context, userID uint64) error {
	key := fmt.Sprintf("user:%d:token", userID)
	err := t.RedisClient.Client.Del(ctx, key).Err()
	if err != nil {
		t.LoggerSugar.Errorw(constants.ErrorToDeleteTokenFromRedis, "error", err.Error(), "key", key)
		return err
	}

	t.LoggerSugar.Infow(constants.SuccessTokenDeleted, "userID", userID)
	return nil
}
