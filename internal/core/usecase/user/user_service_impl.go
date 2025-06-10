package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
)

// Service provides an abstraction for user management, including creating, retrieving, updating, and deleting users, plus authentication handling.
type Service struct {
	userRepository db.UserStore
	tokenService   token.Usecase
	securityHasher security.Store
	logger         logger.Logger
}

// NewUserService creates and returns a new Service instance with the provided dependencies for handling user-related operations.
func NewUserService(userRepo db.UserStore, tokenService token.Usecase, securityHasher security.Store, logger logger.Logger) *Service {
	return &Service{
		userRepository: userRepo,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
	}
}
