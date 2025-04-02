package auth

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type Service struct {
	UserRetriever  db.Retriever
	TokenService   token.Service
	PasswordHasher security.Hasher
	LoggerSugar    *zap.SugaredLogger
	SecretKey      string
}

func NewAuthService(userRetriever db.Retriever, tokenService token.Service, passwordHasher security.Hasher, loggerSugar *zap.SugaredLogger, secretKey string) *Service {
	return &Service{
		UserRetriever:  userRetriever,
		TokenService:   tokenService,
		PasswordHasher: passwordHasher,
		LoggerSugar:    loggerSugar,
		SecretKey:      secretKey,
	}
}

func (s *Service) Login(ctx domain.ContextControl, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error) {
	userDB, err := s.UserRetriever.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.PasswordHasher.ComparePasswords(userDB.Password, passwordReq); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}

	newToken, err := s.TokenService.Create(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain.Token = newToken

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	s.LoggerSugar.Infow(constants.SuccessToLogin, constants.UserID, userDB.ID)
	return userDB, tokenDomain.Token, nil
}

func (s *Service) Logout(ctx domain.ContextControl, token string) error {
	userID, _, err := s.TokenService.Check(ctx, token)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCheckToken, constants.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	if err := s.TokenService.Delete(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToRevokeToken, constants.Error, err.Error(), constants.UserID, userID)
		return err
	}

	s.LoggerSugar.Infow(constants.SuccessUserLoggedOut, constants.UserID, userID)
	return nil
}
