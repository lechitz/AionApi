package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepository  db.UserRepository
	TokenService    token.TokenService
	PasswordService security.PasswordManager
	LoggerSugar     *zap.SugaredLogger
}

func NewUserService(userRepo db.UserRepository, tokenService token.TokenService, passwordService security.PasswordManager, loggerSugar *zap.SugaredLogger) *UserService {
	return &UserService{
		UserRepository:  userRepo,
		TokenService:    tokenService,
		PasswordService: passwordService,
		LoggerSugar:     loggerSugar,
	}
}
