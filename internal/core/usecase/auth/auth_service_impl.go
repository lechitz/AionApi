package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type AuthService struct {
	UserRetriever  db.Retriever
	TokenService   token.Service
	PasswordHasher security.Hasher
	LoggerSugar    *zap.SugaredLogger
	SecretKey      string
}

func NewAuthService(userRetriever db.Retriever, tokenService token.Service, passwordHasher security.Hasher, loggerSugar *zap.SugaredLogger, secretKey string) *AuthService {
	return &AuthService{
		UserRetriever:  userRetriever,
		TokenService:   tokenService,
		PasswordHasher: passwordHasher,
		LoggerSugar:    loggerSugar,
		SecretKey:      secretKey,
	}
}
