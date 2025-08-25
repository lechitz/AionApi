// Package auth provides operations for managing user authentication and token management.
package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service provides authentication operations including login, logout, and user token management.
type Service struct {
	userRepository output.UserRepository
	hasher         output.Hasher
	tokenProvider  output.TokenProvider
	tokenStore     output.TokenStore
	logger         output.ContextLogger
}

// NewService creates and returns a new instance of Service with dependencies for user retrieval, token management, and security operations.
func NewService(userRepository output.UserRepository, tokenStore output.TokenStore, hasher output.Hasher, tokenProvider output.TokenProvider, logger output.ContextLogger) *Service {
	return &Service{
		userRepository: userRepository,
		tokenStore:     tokenStore,
		tokenProvider:  tokenProvider,
		hasher:         hasher,
		logger:         logger,
	}
}
