package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/core/domain"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/ports/output/cache"
	"go.uber.org/zap"
)

type TokenService struct {
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
	TokenRepository cache.ITokenRepository
}

func NewTokenService(tokenRepo cache.ITokenRepository, loggerSugar *zap.SugaredLogger, secretKey string) *TokenService {
	return &TokenService{
		LoggerSugar:     loggerSugar,
		SecretKey:       secretKey,
		TokenRepository: tokenRepo,
	}
}

func (s *TokenService) Create(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	if _, err := s.TokenRepository.Get(ctx, token); err == nil {
		if err := s.TokenRepository.Delete(ctx, token); err != nil {
			s.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := s.generateToken(token.UserID)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToAssignToken, contextkeys.Error, err.Error())
		return "", err
	}

	s.LoggerSugar.Infow(msg.SuccessTokenCreated, contextkeys.UserID, token.UserID)
	return signedToken, nil
}

func (s *TokenService) Check(ctx domain.ContextControl, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		s.LoggerSugar.Errorw(msg.ErrorInvalidToken, contextkeys.Token, token, contextkeys.Error, err)
		return 0, "", fmt.Errorf(msg.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.LoggerSugar.Errorw(msg.ErrorInvalidTokenClaims, contextkeys.Token, token)
		return 0, "", fmt.Errorf(msg.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims[contextkeys.UserID].(float64)
	if !ok {
		s.LoggerSugar.Errorw(msg.ErrorInvalidUserIDClaim, contextkeys.Token, token)
		return 0, "", fmt.Errorf(msg.ErrorInvalidUserIDClaim)
	}
	userID := uint64(userIDFloat)

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	cachedToken, err := s.TokenRepository.Get(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToRetrieveTokenFromCache, contextkeys.Error, err.Error(), contextkeys.UserID, userID)
		return 0, "", fmt.Errorf(msg.ErrorToRetrieveTokenFromCache)
	}

	if cachedToken != token {
		s.LoggerSugar.Errorw(msg.ErrorTokenMismatch, contextkeys.UserID, userID, msg.TokenFromCookie, token, msg.TokenFromCache, cachedToken)
		return 0, "", fmt.Errorf(msg.ErrorTokenMismatch)
	}

	s.LoggerSugar.Infow(msg.SuccessTokenValidated, contextkeys.UserID, userID)
	return userID, cachedToken, nil
}

func (s *TokenService) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Save(ctx, token); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToSaveToken, contextkeys.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(msg.SuccessTokenCreated, contextkeys.UserID, token.UserID)
	return nil
}

func (s *TokenService) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Update(ctx, token); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToUpdateToken, contextkeys.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(msg.SuccessTokenUpdated, contextkeys.UserID, token.UserID)
	return nil
}

func (s *TokenService) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Delete(ctx, token); err != nil {
		s.LoggerSugar.Errorw(msg.ErrorToDeleteToken, contextkeys.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(msg.SuccessTokenDeleted, contextkeys.UserID, token.UserID)
	return nil
}

func (s *TokenService) generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		contextkeys.UserID: userID,
		contextkeys.Exp:    time.Now().Add(contextkeys.ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}
