// Package handler (user) implements HTTP controllers for user-related endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/handlerhelpers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUser handles PUT /user/.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.SpanUpdateUserHandler)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(constants.EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		handlerhelpers.WriteAuthError(ctx, w, span, h.Logger)
		return
	}

	span.AddEvent(constants.EventUserIDFoundInContext,
		trace.WithAttributes(
			attribute.String(tracingkeys.ContextUserID, strconv.FormatUint(userID, 10)),
		),
	)

	var req dto.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handlerhelpers.WriteDecodeError(ctx, w, span, err, h.Logger)
		return
	}

	cmd := req.ToCommand()
	if !cmd.HasUpdates() {
		handlerhelpers.WriteDecodeError(ctx, w, span, constants.ErrNoFieldsToUpdate, h.Logger)
		return
	}

	span.AddEvent(constants.EventUserServiceUpdateUser)

	userUpdated, err := h.UserService.UpdateUser(ctx, userID, cmd)
	if err != nil {
		handlerhelpers.WriteDomainError(ctx, w, span, err, constants.ErrUpdateUser, h.Logger)
		return
	}

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

	span.SetAttributes(attribute.String(commonkeys.UpdatedUsername, userUpdated.Username))
	span.SetStatus(codes.Ok, constants.StatusUserUpdated)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserResponse{
		ID:        userUpdated.ID,
		Name:      &userUpdated.Name,
		Username:  &userUpdated.Username,
		Email:     &userUpdated.Email,
		UpdatedAt: userUpdated.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, constants.MsgUserUpdated)
}
