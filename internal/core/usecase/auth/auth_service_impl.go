package auth

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type AuthService struct {
	UserRetriever  db.UserRetriever
	TokenService   token.TokenService
	SecurityHasher security.SecurityStore
	LoggerSugar    *zap.SugaredLogger
	SecretKey      string
}

func NewAuthService(userRetriever db.UserRetriever, tokenService token.TokenService, securityHasher security.SecurityStore, loggerSugar *zap.SugaredLogger, secretKey string) *AuthService {
	return &AuthService{
		UserRetriever:  userRetriever,
		TokenService:   tokenService,
		SecurityHasher: securityHasher,
		LoggerSugar:    loggerSugar,
		SecretKey:      secretKey,
	}
}
