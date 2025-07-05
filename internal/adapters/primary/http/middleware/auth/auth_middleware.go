// Package auth provides functionality for authentication in HTTP middleware.
package auth

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth/constants"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// MiddlewareAuth provides authentication middleware functionality.
// Combines TokenService for token verification and Logger for logging operations.
type MiddlewareAuth struct {
	tokenService output.TokenRepositoryPort
	logger       output.Logger
	secretKey    string
}

// NewAuthMiddleware creates and initializes middleware for authentication.
func NewAuthMiddleware(
	tokenService output.TokenRepositoryPort,
	logger output.Logger,
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
		tr := otel.Tracer("MiddlewareAuth")
		ctx, span := tr.Start(r.Context(), "Auth")
		defer span.End()

		tokenCookie, err := extractTokenFromCookie(r)
		if err != nil {
			span.SetStatus(codes.Error, "missing token")
			span.SetAttributes(attribute.String("auth.error", err.Error()))
			a.logger.Warnw(constants.ErrorUnauthorizedAccessMissingToken, def.Error, err.Error())
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
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, def.Error, err)
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

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			span.SetStatus(codes.Error, "missing userID")
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)

			return
		}

		userID := uint64(userIDFloat)
		span.SetAttributes(attribute.String("auth.userID", strconv.FormatUint(userID, 10)))

		tokenDomain := entity.TokenDomain{UserID: userID, Token: tokenCookie}

		if _, err := a.tokenService.Get(ctx, tokenDomain); err != nil {
			span.SetStatus(codes.Error, "token not found in cache")
			span.SetAttributes(attribute.String("auth.error", err.Error()))
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, def.Error, err.Error())

			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)

			return
		}

		span.SetStatus(codes.Ok, "authenticated")
		span.SetAttributes(attribute.String("auth.status", "authenticated"))

		newCtx := context.WithValue(ctx, def.CtxUserID, tokenDomain.UserID)
		newCtx = context.WithValue(newCtx, def.CtxToken, tokenCookie)

		a.logger.Infow("auth context: ", newCtx)

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

// extractTokenFromCookie extracts the token from the request cookie.
// Returns the token string or an error if the token is not found or another issue occurs.
func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(def.AuthToken)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
