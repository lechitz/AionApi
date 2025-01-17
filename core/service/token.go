package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/core/domain/entities"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/ports/output/cache"
	"github.com/lechitz/AionApi/ports/output/db"
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

func (service *TokenService) CreateToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) (string, error) {

	_, err := service.TokenRepository.GetToken(ctx, tokenDomain)
	if err == nil {
		if err := service.DeleteToken(ctx, tokenDomain); err != nil {
			service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
			return "", err
		}
	}

	claims := jwt.MapClaims{
		contextkeys.UserID: tokenDomain.UserID,
		"exp":              time.Now().Add(contextkeys.ExpTimeToken).Unix(), // ExpTimeToken é definido em pkg/contextkeys ou similar
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

func (service *TokenService) SaveToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error {
	err := service.TokenRepository.SaveToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToCreateToken, contextkeys.Error, err.Error(), contextkeys.UserID, tokenDomain.UserID)
		return err
	}
	service.LoggerSugar.Infow(msg.SuccessTokenCreated, contextkeys.UserID, tokenDomain.UserID)
	return nil
}

func (service *TokenService) GetToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) (string, error) {
	token, err := service.TokenRepository.GetToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToRetrieveTokenFromRedis, contextkeys.Error, err.Error())
		return "", err
	}

	service.LoggerSugar.Infow(msg.SuccessTokenRetrieved, contextkeys.UserID, tokenDomain.UserID)
	return token, nil
}

func (service *TokenService) UpdateToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error {
	err := service.TokenRepository.UpdateToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToUpdateToken, contextkeys.Error, err.Error())
		return err
	}
	service.LoggerSugar.Infow(msg.SuccessTokenUpdated, contextkeys.UserID, tokenDomain.UserID)
	return nil
}

func (service *TokenService) DeleteToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error {
	err := service.TokenRepository.DeleteToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return err
	}
	service.LoggerSugar.Infow(msg.SuccessTokenDeleted, contextkeys.UserID, tokenDomain.UserID)
	return nil
}

func (service *TokenService) CheckToken(ctx entities.ContextControl, token string) (uint64, string, error) {
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

	tokenDomain := entities.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	tokenFromRedis, err := service.GetToken(ctx, tokenDomain)
	if err != nil {
		service.LoggerSugar.Errorw(msg.ErrorToRetrieveTokenFromRedis, contextkeys.Error, err.Error(), contextkeys.UserID, userID)
		return 0, "", fmt.Errorf(msg.ErrorToRetrieveTokenFromRedis)
	}

	if token != tokenFromRedis {
		service.LoggerSugar.Errorw(msg.ErrorTokenMismatch, contextkeys.UserID, userID, "tokenFromCookie", token, "tokenFromRedis", tokenFromRedis)
		return 0, "", fmt.Errorf(msg.ErrorTokenMismatch)
	}

	service.LoggerSugar.Infow(msg.SuccessTokenValidated, contextkeys.UserID, userID)
	return userID, tokenFromRedis, nil
}
