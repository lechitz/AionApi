package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

type TokenService struct {
	tokenRepository cache.TokenRepositoryPort
	logger          logger.Logger
	configToken     domain.TokenConfig
}

func NewTokenService(tokenRepo cache.TokenRepositoryPort, logger logger.Logger, config domain.TokenConfig) *TokenService {
	return &TokenService{
		tokenRepository: tokenRepo,
		logger:          logger,
		configToken:     config,
	}
}
