// Package user provides operations for managing user creation, retrieval, update, and deletion, as well as authentication and token management.
package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service provides an abstraction for user management, including creating, retrieving, updating, and deleting users, plus authentication handling.
type Service struct {
	userRepository output.UserRepository
	tokenService   input.TokenService
	hashStore      output.Hasher
	logger         output.ContextLogger
}

// NewService creates and returns a new Service instance with the provided dependencies for handling user-related operations.
func NewService(userRepository output.UserRepository, tokenService input.TokenService, hashStore output.Hasher, logger output.ContextLogger) *Service {
	return &Service{
		userRepository: userRepository,
		tokenService:   tokenService,
		hashStore:      hashStore,
		logger:         logger,
	}
}
