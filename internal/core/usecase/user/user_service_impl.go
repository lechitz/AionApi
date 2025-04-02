package user

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepository  db.Repository
	TokenService    token.Service
	PasswordService security.Hasher
	LoggerSugar     *zap.SugaredLogger
}

func NewUserService(userRepo db.Repository, tokenService token.Service, passwordService security.Hasher, loggerSugar *zap.SugaredLogger) *UserService {
	return &UserService{
		UserRepository:  userRepo,
		TokenService:    tokenService,
		PasswordService: passwordService,
		LoggerSugar:     loggerSugar,
	}
}
