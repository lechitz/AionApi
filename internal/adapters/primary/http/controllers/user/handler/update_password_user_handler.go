// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/httputils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUserPassword handles PUT /user/password.
func (h *Handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanUpdatePasswordHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(constants.EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	var req dto.UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	err := handlerhelpers.CheckRequiredFields(map[string]string{
		commonkeys.Password:    req.Password,
		commonkeys.NewPassword: req.NewPassword,
	})
	if err != nil {
		h.Logger.ErrorwCtx(ctx, constants.ErrUpdateUserPasswordValidation,
			commonkeys.Error, err.Error(),
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
		)
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		handlerhelpers.WriteAuthError(ctx, w, span, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Bool(constants.SpanAttrUserPasswordUpdate, true),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	newToken, err := h.UserService.UpdatePassword(ctx, userID, req.Password, req.NewPassword)
	if err != nil {
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrUpdateUser, h.Logger)
		return
	}

	span.AddEvent(constants.EventUserServiceUpdateUserPassword,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	httputils.ClearAuthCookie(w, h.Config.Cookie)

	httputils.SetAuthCookie(w, newToken, h.Config.Cookie)

	span.SetStatus(codes.Ok, constants.StatusUserPasswordUpdated)
	span.SetAttributes(
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
		attribute.Bool(constants.SpanAttrAuthCookieRefreshed, true),
	)

	h.Logger.InfowCtx(ctx, constants.MsgUserPasswordUpdated,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(constants.EventUserPasswordUpdatedSuccess)

	httpresponse.WriteSuccess(w, http.StatusOK, nil, constants.MsgUserPasswordUpdated)
}
