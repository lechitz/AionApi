// Package auth provides operations for managing user authentication and token management.
package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
)

// Service provides authentication operations including login, logout, and user token management.
type Service struct {
	userRetriever  output.UserRetriever
	tokenService   token.Usecase
	securityHasher output.Store
	logger         output.Logger
	secretKey      string
}

// NewAuthService creates and returns a new instance of Service with dependencies for user retrieval, token management, and security operations.
func NewAuthService(userRetriever output.UserStore, tokenService token.Usecase, securityHasher output.Store, logger output.Logger, secretKey string) *Service {
	return &Service{
		userRetriever:  userRetriever,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
		secretKey:      secretKey,
	}
}
