// Package user implements HTTP handlers for user-related endpoints.
package user

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUser handles PUT /user/.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanUpdateUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		handlerhelpers.WriteAuthError(ctx, w, span, h.Logger)
		return
	}

	span.AddEvent(constants.EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	req, err := h.parseUpdateUserRequest(r)
	if err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Username, safeStringPtr(req.Username)),
		attribute.String(commonkeys.Email, safeStringPtr(req.Email)),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	userDomain := buildUserDomainFromUpdate(userID, req)
	span.AddEvent(constants.EventUserServiceUpdateUser)

	userUpdated, err := h.Service.UpdateUser(ctx, userDomain)
	if err != nil {
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrUpdateUser, h.Logger)
		return
	}

	span.SetStatus(codes.Ok, constants.StatusUserUpdated)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	h.Logger.InfowCtx(ctx, constants.MsgUserUpdated,
		commonkeys.UserID, strconv.FormatUint(userUpdated.ID, 10),
		commonkeys.Username, userUpdated.Username,
		commonkeys.Email, userUpdated.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(
		constants.EventUserUpdatedSuccess,
		trace.WithAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userUpdated.ID, 10))),
	)

	h.writeUpdateSuccess(w, span, userUpdated)
}

// safeStringPtr safely dereferences a *string, returning an empty string if nil.
func safeStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
