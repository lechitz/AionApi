// Package auth implements HTTP handlers for authentication endpoints.
package auth

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/httputils"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Logout handles the user logout request: invalidates the token, clears cookies, logs the event, and returns a standard response.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerAuthHandler).
		Start(r.Context(), constants.SpanLogoutHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		err := sharederrors.ErrMissingUserID()
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrMissingUserID)
		h.Logger.ErrorwCtx(ctx, constants.LogMissingUserID, commonkeys.Error, err.Error())
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	tokenVal := ctx.Value(ctxkeys.Token)
	tokenString, ok := tokenVal.(string)
	if !ok || tokenString == "" {
		err := sharederrors.ErrUnauthorized(constants.ErrMissingToken)
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrUnauthorized)
		h.Logger.ErrorwCtx(ctx, constants.LogMissingToken, commonkeys.Error, err.Error())
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	span.AddEvent(
		constants.EventAuthServiceLogout,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	if err := h.Service.Logout(ctx, tokenString); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, constants.LogLogoutFailed, commonkeys.Error, err.Error())
		httpresponse.WriteError(w, err, constants.ErrLogout, h.Logger)
		return
	}

	httputils.ClearAuthCookie(w, h.Config.Cookie)

	tokenPreview := ""
	if len(tokenString) >= 10 {
		tokenPreview = tokenString[:10] + "..."
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Token, tokenPreview),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.SetStatus(codes.Ok, constants.StatusLogoutSuccess)

	h.Logger.InfowCtx(ctx, constants.MsgLogoutSuccess,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		commonkeys.Token, tokenPreview,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(constants.EventLogoutSuccess)

	httpresponse.WriteSuccess(w, http.StatusOK, nil, constants.MsgLogoutSuccess)
}
