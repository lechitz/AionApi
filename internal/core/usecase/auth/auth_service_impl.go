// Package auth provides operations for managing user authentication and token management.
package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service provides authentication operations including login, logout, and user token management.
type Service struct {
	userRetriever  output.UserRetriever
	tokenService   input.TokenService
	securityHasher output.HasherStore
	logger         output.Logger
}

// NewAuthService creates and returns a new instance of Service with dependencies for user retrieval, token management, and security operations.
func NewAuthService(userRetriever output.UserRetriever, tokenService input.TokenService, securityHasher output.HasherStore, logger output.Logger) *Service {
	return &Service{
		userRetriever:  userRetriever,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
	}
}
