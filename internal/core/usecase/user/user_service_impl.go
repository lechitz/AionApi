package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
)

type UserService struct {
	userRepository db.UserStore
	tokenService   token.TokenUsecase
	securityHasher security.SecurityStore
	logger         logger.Logger
}

func NewUserService(userRepo db.UserStore, tokenService token.TokenUsecase, securityHasher security.SecurityStore, logger logger.Logger) *UserService {
	return &UserService{
		userRepository: userRepo,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
	}
}
