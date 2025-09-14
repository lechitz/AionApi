// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/cookies"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDeleteUser handles DELETE /user/.
func (h *Handler) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).
		Start(r.Context(), SpanSoftDeleteUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		helpers.WriteAuthError(ctx, w, span, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(EventUserServiceSoftDeleteUser,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	if err := h.UserService.SoftDeleteUser(ctx, userID); err != nil {
		h.Logger.ErrorwCtx(ctx, ErrSoftDeleteUser,
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err.Error(),
		)
		helpers.WriteDomainError(ctx, w, span, err, ErrSoftDeleteUser, h.Logger)
		return
	}

	cookies.ClearAuthCookie(w, h.Config.Cookie)
	span.SetStatus(codes.Ok, StatusUserSoftDeleted)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusNoContent))

	h.Logger.InfowCtx(ctx, MsgUserSoftDeleted,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventUserSoftDeletedSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		),
	)

	httpresponse.WriteNoContent(w)
}
