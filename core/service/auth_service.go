package service

import (
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/ports/output/db"
	"github.com/lechitz/AionApi/ports/output/security"
	"go.uber.org/zap"
)

type AuthService struct {
	UserRepository  db.IUserRepository
	TokenService    security.ITokenService
	PasswordService security.IPasswordService
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
}

func NewAuthService(
	userRepo db.IUserRepository,
	tokenService security.ITokenService,
	passwordService security.IPasswordService,
	loggerSugar *zap.SugaredLogger,
	secretKey string,
) *AuthService {
	return &AuthService{
		UserRepository:  userRepo,
		TokenService:    tokenService,
		PasswordService: passwordService,
		LoggerSugar:     loggerSugar,
		SecretKey:       secretKey,
	}
}

func (s *AuthService) Login(ctx domain.ContextControl, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error) {

	userDB, err := s.UserRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.PasswordService.ComparePasswords(userDB.Password, passwordReq); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToCompareHashAndPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}

	token, err := s.TokenService.Create(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToCreateToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain.Token = token

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToSaveToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	s.LoggerSugar.Infow(msg.SuccessToLogin, contextkeys.UserID, userDB.ID)
	return userDB, tokenDomain.Token, nil
}

func (s *AuthService) Logout(ctx domain.ContextControl, token string) error {
	userID, _, err := s.TokenService.Check(ctx, token)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToCheckToken, contextkeys.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	if err := s.TokenService.Delete(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorRevokeToken, contextkeys.Error, err.Error(), contextkeys.UserID, userID)
		return err
	}

	s.LoggerSugar.Infow(msg.SuccessUserLoggedOut, contextkeys.UserID, userID)
	return nil
}
