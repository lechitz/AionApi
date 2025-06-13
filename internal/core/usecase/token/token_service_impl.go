package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Service manages token-related operations including creation, validation, and deletion using a repository and logging functionality.
type Service struct {
	tokenRepository cache.TokenRepositoryPort
	logger          logger.Logger
	configToken     domain.TokenConfig
}

// NewTokenService initializes a Service with a token repository, logger, and token configuration for managing token operations.
func NewTokenService(
	tokenRepo cache.TokenRepositoryPort,
	logger logger.Logger,
	config domain.TokenConfig,
) *Service {
	return &Service{
		tokenRepository: tokenRepo,
		logger:          logger,
		configToken:     config,
	}
}
