package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"go.uber.org/zap"
)

type TokenService struct {
	TokenRepository cache.TokenRepositoryPort
	LoggerSugar     *zap.SugaredLogger
	ConfigToken     domain.TokenConfig
}

func NewTokenService(tokenRepo cache.TokenRepositoryPort, logger *zap.SugaredLogger, config domain.TokenConfig) *TokenService {
	return &TokenService{
		TokenRepository: tokenRepo,
		LoggerSugar:     logger,
		ConfigToken:     config,
	}
}
