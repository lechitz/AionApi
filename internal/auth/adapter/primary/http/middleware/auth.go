// Package middleware provides authentication middleware for HTTP.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/auth/core/port/input"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/port/output/logger"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type MiddlewareAuth struct {
	tokenService input.TokenService
	logger       logger.ContextLogger
}

func New(tokenService input.TokenService, logger logger.ContextLogger) *MiddlewareAuth {
	return &MiddlewareAuth{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (a *MiddlewareAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		userID, claims, err := a.tokenService.Validate(ctx, rawToken)
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
