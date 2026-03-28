// Package handler contains HTTP handlers for auth operations.
package handler

import (
	"errors"
	"net/http"

	"github.com/lechitz/aion-api/internal/auth/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Refresh renews the access token using a valid refresh token from cookie.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuthHandler).Start(r.Context(), SpanRefreshHandler)
	defer span.End()

	refreshToken, err := cookies.ExtractRefreshToken(r)
	if err != nil {
		span.RecordError(err)
		h.Logger.ErrorwCtx(ctx, ErrRefresh, "reason", err.Error())
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), h.Logger)
		return
	}
	if refreshToken == "" {
		err := errors.New(ErrMissingRefreshToken)
		span.RecordError(err)
		h.Logger.ErrorwCtx(ctx, ErrRefresh, "reason", err.Error())
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), h.Logger)
		return
	}

	span.SetAttributes(attribute.Bool(AttrRefreshTokenPresent, true))

	span.AddEvent(EventAuthServiceRefresh)
	accessToken, newRefreshToken, err := h.Service.RefreshTokenRenewal(ctx, refreshToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrRefresh)
		h.Logger.ErrorwCtx(ctx, ErrRefresh, "reason", err.Error())
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), h.Logger)
		return
	}

	cookies.SetAuthCookie(w, accessToken, h.Config.Cookie)
	cookies.SetRefreshCookie(w, newRefreshToken, h.Config.Cookie)

	// Return the new access token in the response body so frontend can update in-memory token
	refreshResponse := dto.RefreshResponse{
		Token: accessToken,
	}

	span.AddEvent(EventRefreshSuccess)
	span.SetStatus(codes.Ok, StatusRefreshSuccess)
	h.Logger.InfowCtx(ctx, MsgRefreshSuccess)

	httpresponse.WriteSuccess(w, http.StatusOK, refreshResponse, MsgRefreshSuccess)
}
