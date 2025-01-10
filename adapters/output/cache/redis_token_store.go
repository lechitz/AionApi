package cache

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/infra/cache"
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
		t.LoggerSugar.Errorw("Failed to create token", "error", err.Error())
		return "", err
	}

	key := fmt.Sprintf("user:%d:token", userID)
	err = t.RedisClient.Client.Set(ctx, key, signedToken, 24*time.Hour).Err()
	if err != nil {
		t.LoggerSugar.Errorw("Failed to store token in Redis", "error", err.Error(), "userID", userID)
		return "", err
	}

	t.LoggerSugar.Infow("Token created successfully", "userID", userID, "token", signedToken)
	return signedToken, nil
}

func (t *TokenStore) ValidateToken(ctx context.Context, tokenFromCookie string) (string, uint64, error) {
	parsedToken, err := jwt.Parse(tokenFromCookie, func(token *jwt.Token) (interface{}, error) {
		return t.SecretKey, nil
	})
	if err != nil || !parsedToken.Valid {
		t.LoggerSugar.Errorw("Invalid token", "token", tokenFromCookie, "error", err)
		return "", 0, fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.LoggerSugar.Errorw("Invalid token claims", "token", tokenFromCookie)
		return "", 0, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		t.LoggerSugar.Errorw("Invalid userID claim", "token", tokenFromCookie)
		return "", 0, fmt.Errorf("invalid userID claim")
	}

	userID := uint64(userIDFloat)

	tokenFromRedis, err := t.GetTokenByUserID(ctx, userID)
	if err != nil {
		t.LoggerSugar.Errorw("Error retrieving token from Redis", "error", err.Error(), "userID", userID)
		return "", 0, fmt.Errorf("token not found or revoked")
	}

	if tokenFromCookie != tokenFromRedis {
		t.LoggerSugar.Errorw("Token mismatch in Redis", "userID", userID, "tokenFromCookie", tokenFromCookie, "tokenFromRedis", tokenFromRedis)
		return "", 0, fmt.Errorf("token mismatch")
	}

	return tokenFromCookie, userID, nil
}

func (t *TokenStore) GetTokenByUserID(ctx context.Context, userID uint64) (string, error) {
	key := fmt.Sprintf("user:%d:token", userID)
	value, err := t.RedisClient.Client.Get(ctx, key).Result()
	if err != nil {
		t.LoggerSugar.Errorw("Error to get token from Redis", "error", err.Error(), "key", key)
		return "", err
	}

	t.LoggerSugar.Infow("Token retrieved successfully", "userID", userID, "token", value)
	return value, nil
}

func (t *TokenStore) DeleteTokenByUserID(ctx context.Context, userID uint64) error {
	key := fmt.Sprintf("user:%d:token", userID)
	err := t.RedisClient.Client.Del(ctx, key).Err()
	if err != nil {
		t.LoggerSugar.Errorw("Error to delete token from Redis", "error", err.Error(), "key", key)
		return err
	}

	t.LoggerSugar.Infow("Token deleted successfully", "userID", userID)
	return nil
}
