package token

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// Service manages token-related operations including creation, validation, and deletion using a repository and logging functionality.
type Service struct {
	tokenRepository output.TokenStore
	logger          output.Logger
	secretKey       string
}

// NewTokenService initializes a Service with a token repository, logger, and token configuration for managing token operations.
func NewTokenService(tokenRepo output.TokenStore, logger output.Logger, secretKey config.Secret) *Service {
	return &Service{
		tokenRepository: tokenRepo,
		logger:          logger,
		secretKey:       secretKey.Key,
	}
}
