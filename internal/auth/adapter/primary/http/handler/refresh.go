// Package handler contains HTTP handlers for auth operations.
package handler

import (
	"errors"
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Refresh renews the access token using a valid refresh token from cookie.
// This method is implemented on the `Handler` struct which is defined in
// `0_auth_handler_impl.go` (constructor & struct live there).
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuthHandler).Start(r.Context(), SpanRefreshHandler)
	defer span.End()

	refreshToken, err := cookies.ExtractRefreshToken(r)
	// Ensure we use a non-nil error variable when reporting failures to avoid static analysis warnings
	var authErr error
	if err != nil {
		authErr = err
	} else if refreshToken == "" {
		authErr = errors.New("missing refresh token")
	}
	if authErr != nil {
		// record the original error but surface a standardized unauthorized error to the HTTP layer
		span.RecordError(authErr)
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(authErr.Error()), nil)
		return
	}

	// Do NOT record the raw refresh token or any credentials. Only record presence.
	span.SetAttributes(attribute.Bool("refresh_token_present", true))

	accessToken, newRefreshToken, err := h.Service.RefreshTokenRenewal(ctx, refreshToken)
	if err != nil {
		// record the real error for tracing, but return a 401 Unauthorized to clients.
		span.RecordError(err)
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), nil)
		return
	}

	cookies.SetAuthCookie(w, accessToken, h.Config.Cookie)
	cookies.SetRefreshCookie(w, newRefreshToken, h.Config.Cookie)
	w.WriteHeader(http.StatusOK)
}
