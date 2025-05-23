package auth

import (
	"context"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/auth/constants"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config"
)

type MiddlewareAuth struct {
	tokenService cache.TokenRepositoryPort
	logger       logger.Logger
}

func NewAuthMiddleware(tokenService cache.TokenRepositoryPort, logger logger.Logger) *MiddlewareAuth {
	return &MiddlewareAuth{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			a.logger.Warnw(constants.ErrorUnauthorizedAccessMissingToken, constants.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Setting.Secret.Key), nil
		})
		if err != nil || parsedToken == nil || !parsedToken.Valid {
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, constants.Error, err)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims[constants.UserID].(float64)
		if !ok {
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		userID := uint64(userIDFloat)

		tokenDomain := domain.TokenDomain{
			UserID: userID,
			Token:  tokenCookie,
		}

		_, err = a.tokenService.Get(r.Context(), tokenDomain)
		if err != nil {
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, constants.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), constants.UserID, tokenDomain.UserID)
		newCtx = context.WithValue(newCtx, constants.Token, tokenCookie)

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
