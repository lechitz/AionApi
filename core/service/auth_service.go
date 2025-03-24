package service

import (
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/lechitz/AionApi/ports/output/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository  db.IUserRepository
	TokenRepository cache.ITokenRepository
	TokenService    *TokenService
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
}

func NewAuthService(userRepo db.IUserRepository, tokenRepo cache.ITokenRepository, tokenService *TokenService, loggerSugar *zap.SugaredLogger, secretKey string) *AuthService {
	return &AuthService{
		UserRepository:  userRepo,
		TokenRepository: tokenRepo,
		TokenService:    tokenService,
		LoggerSugar:     loggerSugar,
		SecretKey:       secretKey,
	}
}

func (service *AuthService) Login(ctx domain.ContextControl, userDomain domain.UserDomain, passwordReq string) (domain.UserDomain, string, error) {

	userDB, err := service.UserRepository.GetUserByUsername(ctx, userDomain.Username)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToGetUserByUserName, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err = service.compareHashAndPassword(userDB.Password, passwordReq); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCompareHashAndPassword, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userDB.ID,
	}

	token, err := service.TokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCreateToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain.Token = token

	if err := service.TokenService.SaveToken(ctx, tokenDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToSaveToken, contextkeys.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	return userDB, tokenDomain.Token, nil
}

func (service *AuthService) Logout(ctx domain.ContextControl, token string) error {
	userID, _, err := service.TokenService.CheckToken(ctx, token)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCheckToken, contextkeys.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	if err := service.TokenRepository.DeleteToken(ctx, tokenDomain); err != nil {
		service.LoggerSugar.Errorw(msg.ErrorRevokeToken, contextkeys.Error, err.Error(), contextkeys.UserID, userID)
		return err
	}

	service.LoggerSugar.Infow(msg.SuccessUserLoggedOut, contextkeys.UserID, userID)
	return nil
}

func (service *AuthService) compareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
