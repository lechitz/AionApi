// Package auth provides operations for managing user authentication and token management.
package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service provides authentication operations including login, logout, and user token management.
type Service struct {
	userService    output.UserByUsernameFinder
	tokenService   input.TokenService
	securityHasher output.Hasher
	logger         output.ContextLogger
}

// NewService creates and returns a new instance of Service with dependencies for user retrieval, token management, and security operations.
func NewService(userByUsernameFinder output.UserByUsernameFinder, tokenService input.TokenService, securityHasher output.Hasher, logger output.ContextLogger) *Service {
	return &Service{
		userService:    userByUsernameFinder,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
	}
}
