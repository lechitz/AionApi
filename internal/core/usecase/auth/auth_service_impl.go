package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
)

type AuthService struct {
	userRetriever  db.UserRetriever
	tokenService   token.TokenUsecase
	securityHasher security.SecurityStore
	logger         logger.Logger
	secretKey      string
}

func NewAuthService(userRetriever db.UserStore, tokenService token.TokenUsecase, securityHasher security.SecurityStore, logger logger.Logger, secretKey string) *AuthService {
	return &AuthService{
		userRetriever:  userRetriever,
		tokenService:   tokenService,
		securityHasher: securityHasher,
		logger:         logger,
		secretKey:      secretKey,
	}
}
