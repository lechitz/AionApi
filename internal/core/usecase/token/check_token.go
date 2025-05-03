package token

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

func (s *TokenService) VerifyToken(ctx context.Context, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.configToken.SecretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		s.logger.Errorw(constants.ErrorInvalidToken, constants.Token, token, constants.Error, err)
		return 0, "", fmt.Errorf(constants.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Errorw(constants.ErrorInvalidTokenClaims, constants.Token, token)
		return 0, "", fmt.Errorf(constants.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims[constants.UserID].(float64)
	if !ok {
		s.logger.Errorw(constants.ErrorInvalidUserIDClaim, constants.Token, token)
		return 0, "", fmt.Errorf(constants.ErrorInvalidUserIDClaim)
	}
	userID := uint64(userIDFloat)

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	cachedToken, err := s.tokenRepository.Get(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToRetrieveTokenFromCache, constants.Error, err.Error(), constants.UserID, userID)
		return 0, "", fmt.Errorf(constants.ErrorToRetrieveTokenFromCache)
	}

	if cachedToken != token {
		s.logger.Errorw(constants.ErrorTokenMismatch, constants.UserID, userID, constants.TokenFromCookie, token, constants.TokenFromCache, cachedToken)
		return 0, "", fmt.Errorf(constants.ErrorTokenMismatch)
	}

	s.logger.Infow(constants.SuccessTokenValidated, constants.UserID, userID)
	return userID, cachedToken, nil
}
