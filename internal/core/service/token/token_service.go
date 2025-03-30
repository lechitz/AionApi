package token

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/service/constants"
)

type TokenService interface {
	Create(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
	Check(ctx domain.ContextControl, token string) (uint64, string, error)
	Save(ctx domain.ContextControl, token domain.TokenDomain) error
	Update(ctx domain.ContextControl, token domain.TokenDomain) error
	Delete(ctx domain.ContextControl, token domain.TokenDomain) error
}
type tokenService struct {
	TokenRepository cache.TokenRepository
	LoggerSugar     *zap.SugaredLogger
	SecretKey       string
}

func NewTokenService(
	repo cache.TokenRepository,
	logger *zap.SugaredLogger,
	secretKey string,
) TokenService {
	return &tokenService{
		TokenRepository: repo,
		LoggerSugar:     logger,
		SecretKey:       secretKey,
	}
}

func (s *tokenService) Create(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	if _, err := s.TokenRepository.Get(ctx, token); err == nil {
		if err := s.TokenRepository.Delete(ctx, token); err != nil {
			s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := s.generateToken(token.UserID)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToAssignToken, constants.Error, err.Error())
		return "", err
	}

	s.LoggerSugar.Infow(constants.SuccessTokenCreated, constants.UserID, token.UserID)
	return signedToken, nil
}

func (s *tokenService) Check(ctx domain.ContextControl, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		s.LoggerSugar.Errorw(constants.ErrorInvalidToken, constants.Token, token, constants.Error, err)
		return 0, "", fmt.Errorf(constants.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.LoggerSugar.Errorw(constants.ErrorInvalidTokenClaims, constants.Token, token)
		return 0, "", fmt.Errorf(constants.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims[constants.UserID].(float64)
	if !ok {
		s.LoggerSugar.Errorw(constants.ErrorInvalidUserIDClaim, constants.Token, token)
		return 0, "", fmt.Errorf(constants.ErrorInvalidUserIDClaim)
	}
	userID := uint64(userIDFloat)

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	cachedToken, err := s.TokenRepository.Get(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToRetrieveTokenFromCache, constants.Error, err.Error(), constants.UserID, userID)
		return 0, "", fmt.Errorf(constants.ErrorToRetrieveTokenFromCache)
	}

	if cachedToken != token {
		s.LoggerSugar.Errorw(constants.ErrorTokenMismatch, constants.UserID, userID, constants.TokenFromCookie, token, constants.TokenFromCache, cachedToken)
		return 0, "", fmt.Errorf(constants.ErrorTokenMismatch)
	}

	s.LoggerSugar.Infow(constants.SuccessTokenValidated, constants.UserID, userID)
	return userID, cachedToken, nil
}

func (s *tokenService) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Save(ctx, token); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(constants.SuccessTokenCreated, constants.UserID, token.UserID)
	return nil
}

func (s *tokenService) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Update(ctx, token); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToUpdateToken, constants.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(constants.SuccessTokenUpdated, constants.UserID, token.UserID)
	return nil
}

func (s *tokenService) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Delete(ctx, token); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(constants.SuccessTokenDeleted, constants.UserID, token.UserID)
	return nil
}

func (s *tokenService) generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		constants.UserID: userID,
		constants.Exp:    time.Now().Add(constants.ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}
