// Package handler (auth) implements HTTP controllers for authentication endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Logout invalidates the current session and clears the auth cookie.
//
// @Summary      Logout current user
// @Description  Invalidates the current authenticated session (token or cookie) and clears the auth cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      204  {string}  string  "Logout succeeded (no content)"
// @Header       204  {string}  Set-Cookie  "auth_token=deleted; Path=/; Max-Age=0; HttpOnly; Secure (if enabled)"
// @Failure      401  {string}  string  "Unauthorized or missing user context"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /auth/logout [post].
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuthHandler).
		Start(r.Context(), SpanLogoutHandler)
	defer span.End()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		err := sharederrors.ErrMissingUserID()
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserID)
		h.Logger.ErrorwCtx(ctx, LogMissingUserID, commonkeys.Error, err.Error())
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	span.AddEvent(
		EventAuthServiceLogout,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		),
	)

	if err := h.Service.Logout(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, LogLogoutFailed, commonkeys.Error, err.Error())
		httpresponse.WriteError(w, err, ErrLogout, h.Logger)
		return
	}

	cookies.ClearAuthCookie(w, h.Config.Cookie)
	cookies.ClearRefreshCookie(w, h.Config.Cookie)

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)
	span.SetStatus(codes.Ok, StatusLogoutSuccess)

	h.Logger.InfowCtx(ctx, MsgLogoutSuccess,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
	)
	span.AddEvent(EventLogoutSuccess)

	httpresponse.WriteNoContent(w)
}
