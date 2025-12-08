// Package middleware provides authentication middleware for HTTP.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AuthMiddleware is a middleware that authenticates a user.
type AuthMiddleware struct {
	authService input.AuthService
	logger      logger.ContextLogger
}

// New creates a new instance of AuthMiddleware.
func New(authService input.AuthService, logger logger.ContextLogger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Auth authenticates a user and sets the user ID and token in the context.
func (a *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if already authenticated via service token (S2S)
		if svc, ok := r.Context().Value(ctxkeys.ServiceAccount).(bool); ok && svc {
			if userID, ok := r.Context().Value(ctxkeys.UserID).(uint64); ok {
				a.logger.Infow(MsgS2SAuthBypass,
					commonkeys.UserID, userID,
					commonkeys.URLPath, r.URL.Path,
					commonkeys.Method, r.Method,
				)
			}
			next.ServeHTTP(w, r)
			return
		}

		tr := otel.Tracer(TracerAuthMiddleware)
		ctx, span := tr.Start(r.Context(), SpanAuthMiddleware)
		defer span.End()

		rawToken, err := extractToken(r)
		if err != nil || rawToken == "" {
			span.SetStatus(codes.Error, SpanErrorMissingToken)
			if err == nil {
				err = sharederrors.ErrUnauthorized(ErrorUnauthorizedAccessMissingToken)
			}
			a.logger.WarnwCtx(ctx, ErrorUnauthorizedAccessMissingToken)
			httpresponse.WriteAuthError(w, err, a.logger)
			return
		}

		userID, claims, err := a.authService.Validate(ctx, rawToken)
		if err != nil {
			span.SetStatus(codes.Error, SpanErrorTokenInvalid)
			span.SetAttributes(attribute.String(AttrAuthMiddlewareError, err.Error()))
			a.logger.WarnwCtx(ctx, ErrorUnauthorizedAccessInvalidToken, commonkeys.Error, err.Error())
			httpresponse.WriteAuthError(w, err, a.logger)
			return
		}

		ctx = context.WithValue(ctx, ctxkeys.UserID, userID)
		ctx = context.WithValue(ctx, ctxkeys.Token, rawToken)
		if claims != nil {
			ctx = context.WithValue(ctx, ctxkeys.Claims, claims)
		}

		span.SetStatus(codes.Ok, SpanStatusAuthenticated)
		span.SetAttributes(
			attribute.String(AttrAuthMiddlewareUserID, strconv.FormatUint(userID, 10)),
			attribute.String(AttrAuthMiddlewareStatus, StatusAuthenticated),
		)
		a.logger.InfowCtx(ctx, MsgContextSet, commonkeys.UserID, strconv.FormatUint(userID, 10))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken extracts the token from the request.
func extractToken(r *http.Request) (string, error) {

	// Authorization: Bearer <token>
	if ah := r.Header.Get("Authorization"); ah != "" {
		parts := strings.SplitN(ah, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
			return parts[1], nil
		}
	}

	// Cookie
	if c, err := r.Cookie(commonkeys.AuthTokenCookieName); err == nil && c != nil && c.Value != "" {
		return c.Value, nil
	}

	return "", errors.New(ErrorUnauthorizedAccessMissingToken)
}
