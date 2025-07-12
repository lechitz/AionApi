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
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/httputils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDeleteUser handles DELETE /user/.
func (h *Handler) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanSoftDeleteUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		handlerhelpers.WriteAuthError(ctx, w, span, h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)

	span.AddEvent(constants.EventUserServiceSoftDeleteUser,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	if err := h.Service.SoftDeleteUser(ctx, userID); err != nil {
		h.Logger.ErrorwCtx(ctx, constants.ErrSoftDeleteUser,
			tracingkeys.RequestIPKey, ip,
			tracingkeys.RequestUserAgentKey, userAgent,
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err.Error(),
		)
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrSoftDeleteUser, h.Logger)
		return
	}

	httputils.ClearAuthCookie(w, h.Config.Cookie)
	span.SetStatus(codes.Ok, constants.StatusUserSoftDeleted)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusNoContent))

	h.Logger.InfowCtx(ctx, constants.MsgUserSoftDeleted,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(constants.EventUserSoftDeletedSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		),
	)

	httpresponse.WriteNoContent(w)
}
