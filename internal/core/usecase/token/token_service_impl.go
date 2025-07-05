package token

import (
	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Service manages token-related operations including creation, validation, and deletion using a repository and logging functionality.
type Service struct {
	tokenRepository output.TokenRepositoryPort
	logger          logger.Logger
	configToken     entity.TokenConfig
}

// NewTokenService initializes a Service with a token repository, logger, and token configuration for managing token operations.
func NewTokenService(tokenRepo output.TokenRepositoryPort, logger logger.Logger, config entity.TokenConfig) *Service {
	return &Service{
		tokenRepository: tokenRepo,
		logger:          logger,
		configToken:     config,
	}
}
