package service

import (
	"fmt"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/lechitz/AionApi/ports/output/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
)

type TokenService struct {
	UserRepository  db.IUserRepository
	TokenRepository cache.ITokenRepository
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
}

func NewTokenService(userRepo db.IUserRepository, tokenRepo cache.ITokenRepository, loggerSugar *zap.SugaredLogger, secretKey string) *TokenService {
	return &TokenService{
		UserRepository:  userRepo,
		TokenRepository: tokenRepo,
		LoggerSugar:     loggerSugar,
		SecretKey:       secretKey,
	}
}

func (service *TokenService) CreateToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) (string, error) {

	_, err := service.TokenRepository.GetToken(ctx, tokenDomain)
	if err == nil {
		if err := service.DeleteToken(ctx, tokenDomain); err != nil {
			service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
			return "", err
		}
	}

	claims := jwt.MapClaims{
		contextkeys.UserID: tokenDomain.UserID,
		contextkeys.Exp:    time.Now().Add(contextkeys.ExpTimeToken).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := newToken.SignedString([]byte(service.SecretKey))
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToAssignToken, contextkeys.Error, err.Error())
		return "", err
	}

	service.LoggerSugar.Infow(msg.SuccessTokenCreated, contextkeys.UserID, tokenDomain.UserID)
	return signedToken, nil
}

func (service *TokenService) SaveToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error {
	err := service.TokenRepository.SaveToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCreateToken, contextkeys.Error, err.Error(), contextkeys.UserID, tokenDomain.UserID)
		return err
	}
	service.LoggerSugar.Infow(msg.SuccessTokenCreated, contextkeys.UserID, tokenDomain.UserID)
	return nil
}

func (service *TokenService) GetToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) (string, error) {
	token, err := service.TokenRepository.GetToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToRetrieveTokenFromCache, contextkeys.Error, err.Error())
		return "", err
	}

	service.LoggerSugar.Infow(msg.SuccessTokenRetrieved, contextkeys.UserID, tokenDomain.UserID)
	return token, nil
}

func (service *TokenService) UpdateToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error {
	err := service.TokenRepository.UpdateToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToUpdateToken, contextkeys.Error, err.Error())
		return err
	}
	service.LoggerSugar.Infow(msg.SuccessTokenUpdated, contextkeys.UserID, tokenDomain.UserID)
	return nil
}

func (service *TokenService) DeleteToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error {
	err := service.TokenRepository.DeleteToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return err
	}
	service.LoggerSugar.Infow(msg.SuccessTokenDeleted, contextkeys.UserID, tokenDomain.UserID)
	return nil
}

func (service *TokenService) CheckToken(ctx domain.ContextControl, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		return []byte(service.SecretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		service.LoggerSugar.Errorw(msg.ErrorInvalidToken, contextkeys.Token, token, contextkeys.Error, err)
		return 0, "", fmt.Errorf(msg.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		service.LoggerSugar.Errorw(msg.ErrorInvalidTokenClaims, contextkeys.Token, token)
		return 0, "", fmt.Errorf(msg.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims[contextkeys.UserID].(float64)
	if !ok {
		service.LoggerSugar.Errorw(msg.ErrorInvalidUserIDClaim, contextkeys.Token, token)
		return 0, "", fmt.Errorf(msg.ErrorInvalidUserIDClaim)
	}

	userID := uint64(userIDFloat)

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	tokenFromCache, err := service.GetToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToRetrieveTokenFromCache, contextkeys.Error, err.Error(), contextkeys.UserID, userID)
		return 0, "", fmt.Errorf(msg.ErrorToRetrieveTokenFromCache)
	}

	if token != tokenFromCache {
		service.LoggerSugar.Errorw(msg.ErrorTokenMismatch, contextkeys.UserID, userID, msg.TokenFromCookie, token, msg.TokenFromCache, tokenFromCache)
		return 0, "", fmt.Errorf(msg.ErrorTokenMismatch)
	}

	service.LoggerSugar.Infow(msg.SuccessTokenValidated, contextkeys.UserID, userID)
	return userID, tokenFromCache, nil
}
