// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/primary/http/dto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUserPassword updates the authenticated user's password and refreshes the auth cookie.
//
// @Summary      Update current user's password
// @Description  Validates the current password and updates it to a new one. On success, the auth cookie is refreshed.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body  body      dto.UpdatePasswordUserRequest  true  "Current and new passwords"
// @Success      200   {string}  string                         "Password updated; auth cookie refreshed"
// @Header       200   {string}  Set-Cookie                     "auth_token=<new>; Path=/; HttpOnly; Secure (if enabled)"
// @Failure      400   {string}  string                         "Invalid request payload or validation error"
// @Failure      401   {string}  string                         "Unauthorized or missing user context"
// @Failure      409   {string}  string                         "Conflict updating password"
// @Failure      500   {string}  string                         "Internal server error"
// @Router       /user/password [put].
func (h *Handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).
		Start(r.Context(), SpanUpdatePasswordHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	var req dto.UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	err := CheckRequiredFields(map[string]string{
		commonkeys.Password:    req.Password,
		commonkeys.NewPassword: req.NewPassword,
	})
	if err != nil {
		h.Logger.ErrorwCtx(ctx, ErrUpdateUserPasswordValidation,
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		err := sharederrors.ErrMissingUserID()
		httpresponse.WriteAuthErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Bool(SpanAttrUserPasswordUpdate, true),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	newToken, err := h.UserService.UpdatePassword(ctx, userID, req.Password, req.NewPassword)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrUpdateUser, h.Logger)
		return
	}

	span.AddEvent(EventUserServiceUpdateUserPassword,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	cookies.ClearAuthCookie(w, h.Config.Cookie)
	cookies.SetAuthCookie(w, newToken, h.Config.Cookie)

	span.SetStatus(codes.Ok, StatusUserPasswordUpdated)
	span.SetAttributes(
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
		attribute.Bool(SpanAttrAuthCookieRefreshed, true),
	)

	h.Logger.InfowCtx(ctx, MsgUserPasswordUpdated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventUserPasswordUpdatedSuccess)

	httpresponse.WriteSuccess(w, http.StatusOK, nil, MsgUserPasswordUpdated)
}
