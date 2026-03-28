// Package handler is the handler for the admin context in the application.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/aion-api/internal/admin/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateUserRoles updates the roles for a specific user.
//
// @Summary      Update user roles (Legacy)
// @Description  Generic role update endpoint with minimal hierarchy validation.
// @Description  **DEPRECATED**: Prefer specific endpoints for better validation:
// @Description  - PUT /admin/users/{id}/promote-admin (adds admin role)
// @Description  - PUT /admin/users/{id}/demote-admin (removes admin role)
// @Description  - PUT /admin/users/{id}/block (blocks user)
// @Description  - PUT /admin/users/{id}/unblock (unblocks user)
// @Description  This endpoint provides maximum flexibility but less validation.
// @Description  Use only for special administrative operations or bulk scripts.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        user_id  path      integer                      true  "User ID (uint64)"
// @Param        body     body      dto.UpdateUserRolesRequest   true  "New roles to assign"
// @Success      200      {object}  dto.UpdateUserRolesResponse  "Roles updated"
// @Failure      400      {string}  string                       "Invalid request or user_id"
// @Failure      401      {string}  string                       "Unauthorized"
// @Failure      403      {string}  string                       "Forbidden - not an admin"
// @Failure      404      {string}  string                       "User not found"
// @Failure      500      {string}  string                       "Internal server error"
// @Deprecated
// @Router       /admin/users/{user_id}/roles [put].
func (h *Handler) UpdateUserRoles(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerName).
		Start(r.Context(), SpanUpdateUserRoles)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	userIDParam := chi.URLParam(r, commonkeys.UserID)
	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, ErrMissingUserIDParam)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam, commonkeys.Error, err.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, ErrInvalidUserIDParam)
		span.RecordError(validationErr)
		span.SetStatus(codes.Error, ErrInvalidUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, validationErr, h.Logger)
		return
	}

	span.AddEvent(EventUserIDExtracted,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		),
	)

	var req dto.UpdateUserRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	if err := req.Validate(); err != nil {
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	cmd := req.ToCommand(userID)

	span.AddEvent(EventAdminServiceUpdateRoles)
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.StringSlice(commonkeys.Roles, cmd.Roles),
	)

	updatedUser, err := h.AdminService.UpdateUserRoles(ctx, cmd)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrUpdateUserRoles, h.Logger)
		return
	}

	h.Logger.InfowCtx(ctx, MsgRolesUpdated,
		commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10),
		commonkeys.Username, updatedUser.Username,
		commonkeys.Roles, updatedUser.Roles,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventRolesUpdatedSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
			attribute.StringSlice(commonkeys.Roles, updatedUser.Roles),
		),
	)

	span.SetStatus(codes.Ok, StatusRolesUpdated)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserRolesResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Roles:     updatedUser.Roles,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgRolesUpdated)
}
