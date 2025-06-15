// Package auth provides functionality for authentication in HTTP middleware.
package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// MiddlewareAuth provides authentication middleware functionality.
// Combines TokenService for token verification and Logger for logging operations.
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
		tr := otel.Tracer("AionApi/Middleware")
		ctx, span := tr.Start(r.Context(), "AuthMiddleware")
		defer span.End()

		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			span.SetStatus(codes.Error, "missing token")
			span.SetAttributes(attribute.String("auth.error", err.Error()))
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
			span.SetStatus(codes.Error, "invalid token")
			if err != nil {
				span.SetAttributes(attribute.String("auth.error", err.Error()))
			}
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, constants.Error, err)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			span.SetStatus(codes.Error, "invalid claims")
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims[constants.UserIDKey].(float64)
		if !ok {
			span.SetStatus(codes.Error, "missing userID")
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}
		userID := uint64(userIDFloat)
		span.SetAttributes(attribute.String("auth.userID", strconv.FormatUint(userID, 10)))

		tokenDomain := domain.TokenDomain{UserID: userID, Token: tokenCookie}

		if _, err := a.tokenService.Get(ctx, tokenDomain); err != nil {
			span.SetStatus(codes.Error, "token not found in cache")
			span.SetAttributes(attribute.String("auth.error", err.Error()))
			a.logger.Warnw(
				constants.ErrorUnauthorizedAccessInvalidToken,
				constants.Error,
				err.Error(),
			)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		span.SetStatus(codes.Ok, "authenticated")
		span.SetAttributes(attribute.String("auth.status", "authenticated"))

		newCtx := context.WithValue(ctx, constants.UserIDCtxKey, tokenDomain.UserID)
		newCtx = context.WithValue(newCtx, constants.TokenCtxKey, tokenCookie)

		a.logger.Infow("auth context: ", newCtx)

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
