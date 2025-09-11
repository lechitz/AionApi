// Package usecase (user) provides operations for managing user creation, retrieval, update, and deletion, as well as authentication and token management.
package usecase

import (
	authOutput "github.com/lechitz/AionApi/internal/auth/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/hasher"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	userOutput "github.com/lechitz/AionApi/internal/user/core/ports/output"
)

// Service provides an abstraction for user management, including creating, retrieving, updating, and deleting users, plus authentication handling.
type Service struct {
	userRepository userOutput.UserRepository
	authStore      authOutput.AuthStore
	tokenProvider  authOutput.AuthProvider
	hasher         hasher.Hasher
	logger         logger.ContextLogger
}

// NewService creates and returns a new Service instance with the provided dependencies for handling user-related operations.
func NewService(userRepository userOutput.UserRepository, authStore authOutput.AuthStore, tokenProvider authOutput.AuthProvider, hasher hasher.Hasher, logger logger.ContextLogger) *Service {
	return &Service{
		userRepository: userRepository,
		authStore:      authStore,
		tokenProvider:  tokenProvider,
		hasher:         hasher,
		logger:         logger,
	}
}
