package token

import (
	"github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	"go.uber.org/zap"
)

type TokenService struct {
	TokenRepository cache.TokenRepository
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
}

func NewTokenService(tokenRepo cache.TokenRepository, logger *zap.SugaredLogger, secretKey string) *TokenService {
	return &TokenService{
		TokenRepository: tokenRepo,
		LoggerSugar:     logger,
		SecretKey:       secretKey,
	}
}
