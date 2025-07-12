// Package token defines use cases for token-related operations.
package token

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// GetToken validates the provided token, checks for a match in the cache, and returns the associated user ID and token or an error.
func (s *Service) GetToken(ctx context.Context, token string) (uint64, string, error) {
	parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		s.logger.Errorw(constants.ErrorInvalidToken, commonkeys.Token, token[0:10]+"...", commonkeys.Error, err)
		return 0, "", errors.New(constants.ErrorInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Errorw(constants.ErrorInvalidTokenClaims, commonkeys.Token, token)
		return 0, "", errors.New(constants.ErrorInvalidTokenClaims)
	}

	userIDFloat, ok := claims[claimskeys.UserID].(float64)
	if !ok {
		s.logger.Errorw(constants.ErrorInvalidUserIDClaim, commonkeys.Token, token)
		return 0, "", errors.New(constants.ErrorInvalidUserIDClaim)
	}
	userID := uint64(userIDFloat)

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	cachedToken, err := s.tokenRepository.Get(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToRetrieveTokenFromCache, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return 0, "", errors.New(constants.ErrorToRetrieveTokenFromCache)
	}

	if cachedToken != token {
		s.logger.Errorw(
			constants.ErrorTokenMismatch,
			commonkeys.UserID,
			strconv.FormatUint(userID, 10),
			commonkeys.TokenFromCookie,
			token,
			commonkeys.TokenFromCache,
			cachedToken,
		)
		return 0, "", errors.New(constants.ErrorTokenMismatch)
	}

	s.logger.Infow(constants.SuccessTokenValidated, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return userID, cachedToken, nil
}
