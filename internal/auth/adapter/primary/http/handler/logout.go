// Package handler (auth) implements HTTP controllers for authentication endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/cookies"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Logout handles the user logout request: invalidates the token, clears cookies, logs the event, and returns a standard response.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuthHandler).
		Start(r.Context(), SpanLogoutHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		err := sharederrors.ErrMissingUserID()
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserID)
		h.Logger.ErrorwCtx(ctx, LogMissingUserID, commonkeys.Error, err.Error())
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	tokenPreview := ""
	if tokenVal, ok := ctx.Value(ctxkeys.Token).(string); ok && len(tokenVal) >= 10 {
		tokenPreview = tokenVal[:10] + "..."
	}

	span.AddEvent(
		EventAuthServiceLogout,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
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

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Token, tokenPreview),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.SetStatus(codes.Ok, StatusLogoutSuccess)

	h.Logger.InfowCtx(ctx, MsgLogoutSuccess,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.Token, tokenPreview,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)
	span.AddEvent(EventLogoutSuccess)

	httpresponse.WriteSuccess(w, http.StatusOK, nil, MsgLogoutSuccess)
}
