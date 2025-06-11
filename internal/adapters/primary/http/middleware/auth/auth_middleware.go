// Package auth provides functionality for authentication in HTTP middleware.
package auth

import (
	"context"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth/constants"
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/contextutil"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// MiddlewareAuth provides functionality for authentication in HTTP middleware.
type MiddlewareAuth struct {
	tokenService cache.TokenRepositoryPort
	logger       logger.Logger
	secretKey    string
}

// NewAuthMiddleware creates and initializes middleware for authentication.
func NewAuthMiddleware(
	tokenService cache.TokenRepositoryPort,
	logger logger.Logger,
	secretKey string,
) *MiddlewareAuth {
	return &MiddlewareAuth{
		tokenService: tokenService,
		logger:       logger,
		secretKey:    secretKey,
	}
}

// Auth validates JWT tokens and attaches user context.
func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			a.logger.Warnw(
				constants.ErrorUnauthorizedAccessMissingToken,
				constants.Error,
				err.Error(),
			)
			http.Error(w, constants.ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.Parse(tokenCookie, func(_ *jwt.Token) (interface{}, error) {
			return []byte(a.secretKey), nil
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

		tokenDomain := domain.TokenDomain{UserID: userID, Token: tokenCookie}

		if _, err := a.tokenService.Get(r.Context(), tokenDomain); err != nil {
			a.logger.Warnw(
				constants.ErrorUnauthorizedAccessInvalidToken,
				constants.Error,
				err.Error(),
			)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), contextutil.UserIDKey, tokenDomain.UserID)

		ctx := context.WithValue(newCtx, contextutil.TokenKey, tokenCookie)

		a.logger.Infow("auth context: ", ctx)

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

// extractTokenFromCookie extracts the token from the request cookie.
// Returns the token string or an error if the token is not found or another issue occurs.
func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
