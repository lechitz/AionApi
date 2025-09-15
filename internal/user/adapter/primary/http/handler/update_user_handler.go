// Package handler (user) implements HTTP controllers for user-related endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUser updates the current authenticated user's profile.
//
// @Summary      Update current user
// @Description  Updates the authenticated user's profile fields. At least one field must be provided.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body  body      dto.UpdateUserRequest   true  "Fields to update"
// @Success      200   {object}  dto.UpdateUserResponse        "User updated"
// @Failure      400   {string}  string                       "Invalid request or no fields to update"
// @Failure      401   {string}  string                       "Unauthorized or missing user context"
// @Failure      404   {string}  string                       "User not found"
// @Failure      409   {string}  string                       "Conflict updating user"
// @Failure      500   {string}  string                       "Internal server error"
// @Router       /user [put].
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
		err := sharederrors.ErrMissingUserID()
		httpresponse.WriteAuthErrorSpan(ctx, w, span, err, h.Logger)
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
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	cmd := req.ToCommand()
	if !cmd.HasUpdates() {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, ErrNoFieldsToUpdate, h.Logger)
		return
	}

	span.AddEvent(EventUserServiceUpdateUser)

	userUpdated, err := h.UserService.UpdateUser(ctx, userID, cmd)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrUpdateUser, h.Logger)
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
