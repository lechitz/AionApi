// Package handler (user) implements HTTP controllers for user-related endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	"github.com/lechitz/AionApi/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUser handles PUT /user/.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).
		Start(r.Context(), SpanUpdateUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		helpers.WriteAuthError(ctx, w, span, h.Logger)
		return
	}

	span.AddEvent(EventUserIDFoundInContext,
		trace.WithAttributes(
			attribute.String(tracingkeys.ContextUserID, strconv.FormatUint(userID, 10)),
		),
	)

	var req dto.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	cmd := req.ToCommand()
	if !cmd.HasUpdates() {
		helpers.WriteDecodeError(ctx, w, span, ErrNoFieldsToUpdate, h.Logger)
		return
	}

	span.AddEvent(EventUserServiceUpdateUser)

	userUpdated, err := h.UserService.UpdateUser(ctx, userID, cmd)
	if err != nil {
		helpers.WriteDomainError(ctx, w, span, err, ErrUpdateUser, h.Logger)
		return
	}

	h.Logger.InfowCtx(ctx, MsgUserUpdated,
		commonkeys.UserID, strconv.FormatUint(userUpdated.ID, 10),
		commonkeys.Username, userUpdated.Username,
		commonkeys.Email, userUpdated.Email,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(
		EventUserUpdatedSuccess,
		trace.WithAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userUpdated.ID, 10))),
	)

	span.SetAttributes(attribute.String(commonkeys.UpdatedUsername, userUpdated.Username))
	span.SetStatus(codes.Ok, StatusUserUpdated)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserResponse{
		ID:        userUpdated.ID,
		Name:      &userUpdated.Name,
		Username:  &userUpdated.Username,
		Email:     &userUpdated.Email,
		UpdatedAt: userUpdated.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUserUpdated)
}
