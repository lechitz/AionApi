// Package authmiddleware provides functionality for authentication in HTTP middleware.
package authmiddleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"

	"github.com/lechitz/AionApi/internal/shared/contextutils"

	"github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/authmiddleware/constants"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// MiddlewareAuth provides authentication middleware functionality.
type MiddlewareAuth struct {
	tokenService         output.TokenStore
	logger               output.ContextLogger
	tokenClaimsExtractor output.TokenClaimsExtractor
}

// New creates and initializes middleware for authentication.
func New(tokenService output.TokenStore, logger output.ContextLogger, tokenClaimsExtractor output.TokenClaimsExtractor) *MiddlewareAuth {
	return &MiddlewareAuth{
		tokenService:         tokenService,
		logger:               logger,
		tokenClaimsExtractor: tokenClaimsExtractor,
	}
}

// Auth validates JWT tokens and attaches user context.
func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := otel.Tracer("MiddlewareAuth")
		ctx, span := tr.Start(r.Context(), "Auth")
		defer span.End()

		claims, err := a.tokenClaimsExtractor.ExtractFromRequest(r)
		if err != nil {
			span.SetStatus(codes.Error, "missing or invalid token")
			span.SetAttributes(attribute.String("authmiddleware.error", err.Error()))
			a.logger.Warnw(constants.ErrorUnauthorizedAccessMissingToken, commonkeys.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
			return
		}

		newCtx := contextutils.InjectUserIntoContext(ctx, claims)

		userID, ok := newCtx.Value(ctxkeys.UserID).(uint64)
		if !ok {
			span.SetStatus(codes.Error, "missing userID in claims")
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken)
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		span.SetAttributes(attribute.String("authmiddleware.userID", strconv.FormatUint(userID, 10)))

		tokenVal, ok := ctx.Value(ctxkeys.Token).(string)
		if !ok || tokenVal == "" {
			tokenVal, err = security.ExtractTokenFromCookie(r)
			if err != nil {
				span.SetStatus(codes.Error, "token missing in context and cookie")
				span.SetAttributes(attribute.String("authmiddleware.error", err.Error()))
				a.logger.Warnw(constants.ErrorUnauthorizedAccessMissingToken, commonkeys.Error, err.Error())
				http.Error(w, constants.ErrorUnauthorizedAccessMissingToken, http.StatusUnauthorized)
				return
			}
		}

		tokenDomain := domain.TokenDomain{UserID: userID, Token: tokenVal}

		if _, err := a.tokenService.Get(ctx, tokenDomain); err != nil {
			span.SetStatus(codes.Error, "token not found in cache")
			span.SetAttributes(attribute.String("authmiddleware.error", err.Error()))
			a.logger.Warnw(constants.ErrorUnauthorizedAccessInvalidToken, commonkeys.Error, err.Error())
			http.Error(w, constants.ErrorUnauthorizedAccessInvalidToken, http.StatusUnauthorized)
			return
		}

		span.SetStatus(codes.Ok, "authenticated")
		span.SetAttributes(attribute.String("authmiddleware.status", "authenticated"))

		newCtx = context.WithValue(newCtx, ctxkeys.Token, tokenDomain.Token)

		a.logger.Infow("authmiddleware context set", commonkeys.UserID, strconv.FormatUint(userID, 10))

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
