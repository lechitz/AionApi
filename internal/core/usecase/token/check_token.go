// Package token defines use cases for token-related operations.
package token

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// VerifyToken validates the provided token, checks for a match in the cache, and returns the associated user ID and token or an error.
func (s *Service) VerifyToken(ctx context.Context, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.configToken.SecretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		s.logger.Errorw(constants.ErrorInvalidToken, def.CtxToken, token, def.Error, err)
		return 0, "", errors.New(constants.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Errorw(constants.ErrorInvalidTokenClaims, def.CtxToken, token)
		return 0, "", errors.New(constants.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims[constants.UserID].(float64)
	if !ok {
		s.logger.Errorw(constants.ErrorInvalidUserIDClaim, def.CtxToken, token)
		return 0, "", errors.New(constants.ErrorInvalidUserIDClaim)
	}
	userID := uint64(userIDFloat)

	tokenDomain := entity.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	cachedToken, err := s.tokenRepository.Get(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToRetrieveTokenFromCache, def.Error, err.Error(), def.CtxUserID, userID)
		return 0, "", errors.New(constants.ErrorToRetrieveTokenFromCache)
	}

	if cachedToken != token {
		s.logger.Errorw(constants.ErrorTokenMismatch, def.CtxUserID, userID, constants.TokenFromCookie, token, constants.TokenFromCache, cachedToken) // TODO: AVALIAR !
		return 0, "", errors.New(constants.ErrorTokenMismatch)
	}

	s.logger.Infow(constants.SuccessTokenValidated, def.CtxUserID, userID)
	return userID, cachedToken, nil
}
