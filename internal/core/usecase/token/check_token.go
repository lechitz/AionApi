package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

func (s *TokenService) VerifyToken(ctx domain.ContextControl, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.ConfigToken.SecretKey), nil
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
