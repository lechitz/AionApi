package auth

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	dbports "github.com/lechitz/AionApi/internal/core/ports/output/db"
	securityports "github.com/lechitz/AionApi/internal/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"go.uber.org/zap"
)

type AuthService struct {
	UserRepository  dbports.UserRepository
	TokenService    token.TokenService
	PasswordService securityports.PasswordManager
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
}

func NewAuthService(
	userRepo dbports.UserRepository,
	tokenService token.TokenService,
	passwordService securityports.PasswordManager,
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
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.PasswordService.ComparePasswords(userDB.Password, passwordReq); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}
	token, err := s.TokenService.Create(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}
	tokenDomain.Token = token

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	s.LoggerSugar.Infow(constants.SuccessToLogin, constants.UserID, userDB.ID)
	return userDB, tokenDomain.Token, nil
}

func (s *AuthService) Logout(ctx domain.ContextControl, token string) error {
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
