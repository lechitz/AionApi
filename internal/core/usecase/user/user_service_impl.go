// Package user provides operations for managing user creation, retrieval, update, and deletion, as well as authentication and token management.
package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
)

// Service provides an abstraction for user management, including creating, retrieving, updating, and deleting users, plus authentication handling.
type Service struct {
	userRepository output.UserStore
	tokenService   token.Usecase
	securityHasher security.Store
	logger         logger.Logger
}

// NewUserService creates and returns a new Service instance with the provided dependencies for handling user-related operations.
func NewUserService(userRepo output.UserStore, tokenService token.Usecase, securityHasher security.Store, logger logger.Logger) *Service {
	return &Service{
		userRepository: userRepo,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
	}
}
