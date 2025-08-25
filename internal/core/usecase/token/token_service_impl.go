package token

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service manages token-related operations including creation, validation, and deletion using a repository and logging functionality.
type Service struct {
	tokenStore    output.TokenStore
	tokenProvider output.TokenProvider
	logger        output.ContextLogger
}

// NewService initializes a Service with a token repository, provider, and logger.
func NewService(tokenStore output.TokenStore, tokenProvider output.TokenProvider, logger output.ContextLogger) *Service {
	return &Service{
		tokenStore:    tokenStore,
		tokenProvider: tokenProvider,
		logger:        logger,
	}
}
