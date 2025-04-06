package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepository db.UserRepository
	TokenService   token.TokenService
	SecurityHasher security.SecurityStore
	LoggerSugar    *zap.SugaredLogger
}

func NewUserService(userRepo db.UserRepository, tokenService token.TokenService, securityHasher security.SecurityStore, loggerSugar *zap.SugaredLogger) *UserService {
	return &UserService{
		UserRepository: userRepo,
		TokenService:   tokenService,
		SecurityHasher: securityHasher,
		LoggerSugar:    loggerSugar,
	}
}
