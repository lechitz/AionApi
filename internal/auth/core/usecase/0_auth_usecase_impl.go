// Package usecase (auth) provides operations for managing user authentication and token management.
package usecase

import (
	authOutput "github.com/lechitz/AionApi/internal/auth/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/hasher"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	userOutput "github.com/lechitz/AionApi/internal/user/core/ports/output"
)

// Service provides authentication operations including login, logout, and user token management.
type Service struct {
	userRepository userOutput.UserRepository
	authStore      authOutput.AuthStore
	authProvider   authOutput.AuthProvider
	hasher         hasher.Hasher
	logger         logger.ContextLogger
}

// NewService creates and returns a new instance of Service with dependencies for user retrieval, token management, and security operations.
func NewService(userRepository userOutput.UserRepository, authStore authOutput.AuthStore, authProvider authOutput.AuthProvider, hasher hasher.Hasher, logger logger.ContextLogger) *Service {
	return &Service{
		userRepository: userRepository,
		authStore:      authStore,
		authProvider:   authProvider,
		hasher:         hasher,
		logger:         logger,
	}
}
