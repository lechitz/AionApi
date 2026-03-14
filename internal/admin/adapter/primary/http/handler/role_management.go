// Package handler is the handler for the admin context in the application.
package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/admin/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/admin/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// PromoteToAdmin promotes a user to admin role.
//
// @Summary      Promote user to admin
// @Description  Adds admin role to a user. Only admins and owners can promote to admin.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        user_id  path      integer                      true  "User ID (uint64)"
// @Success      200      {object}  dto.UpdateUserRolesResponse  "User promoted to admin"
// @Failure      400      {string}  string                       "Invalid user_id"
// @Failure      401      {string}  string                       "Unauthorized"
// @Failure      403      {string}  string                       "Forbidden - insufficient privileges"
// @Failure      404      {string}  string                       "User not found"
// @Failure      500      {string}  string                       "Internal server error"
// @Router       /admin/users/{user_id}/promote-admin [put]
//
//nolint:dupl // role management handlers follow similar pattern.
func (h *Handler) PromoteToAdmin(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerName).
		Start(r.Context(), SpanPromoteToAdmin)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	// Extract target user ID from path
	userIDParam := chi.URLParam(r, commonkeys.UserID)
	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, ErrMissingUserIDParam)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam, commonkeys.Error, err.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, ErrInvalidUserIDParam)
		span.RecordError(validationErr)
		span.SetStatus(codes.Error, ErrInvalidUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, validationErr, h.Logger)
		return
	}

	// Extract actor user ID and roles from context (set by auth middleware)
	actorUserID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok {
		err := sharederrors.ErrUnauthorized(ErrMissingActorUserID)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingActorUserID)
		h.Logger.ErrorwCtx(ctx, ErrMissingActorUserID)
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	actorRoles := extractRolesFromContext(ctx)

	span.AddEvent(EventUserIDExtracted,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
			attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
			attribute.StringSlice(commonkeys.Roles, actorRoles),
		),
	)

	cmd := input.PromoteToAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  actorRoles,
	}

	span.AddEvent(EventAdminServicePromoteToAdmin)
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
	)

	// Call service
	updatedUser, err := h.AdminService.PromoteToAdmin(ctx, cmd)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrPromoteToAdminFailed, h.Logger)
		return
	}

	h.Logger.InfowCtx(ctx, MsgUserPromotedToAdmin,
		commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10),
		commonkeys.Username, updatedUser.Username,
		commonkeys.Roles, updatedUser.Roles,
		LogKeyActorUserID, actorUserID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventPromoteToAdminSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
			attribute.StringSlice(commonkeys.Roles, updatedUser.Roles),
		),
	)

	span.SetStatus(codes.Ok, StatusPromotedToAdmin)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserRolesResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Roles:     updatedUser.Roles,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUserPromotedToAdmin)
}

// DemoteFromAdmin removes admin role from a user.
//
// @Summary      Demote user from admin
// @Description  Removes admin role from a user. Only admins and owners can demote.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        user_id  path      integer                      true  "User ID (uint64)"
// @Success      200      {object}  dto.UpdateUserRolesResponse  "User demoted from admin"
// @Failure      400      {string}  string                       "Invalid user_id"
// @Failure      401      {string}  string                       "Unauthorized"
// @Failure      403      {string}  string                       "Forbidden - insufficient privileges"
// @Failure      404      {string}  string                       "User not found"
// @Failure      500      {string}  string                       "Internal server error"
// @Router       /admin/users/{user_id}/demote-admin [put]
//
//nolint:dupl // role management handlers follow similar pattern.
func (h *Handler) DemoteFromAdmin(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerName).
		Start(r.Context(), SpanDemoteFromAdmin)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	// Extract target user ID from path
	userIDParam := chi.URLParam(r, commonkeys.UserID)
	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, ErrMissingUserIDParam)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam, commonkeys.Error, err.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, ErrInvalidUserIDParam)
		span.RecordError(validationErr)
		span.SetStatus(codes.Error, ErrInvalidUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, validationErr, h.Logger)
		return
	}

	// Extract actor user ID and roles from context
	actorUserID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok {
		err := sharederrors.ErrUnauthorized(ErrMissingActorUserID)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingActorUserID)
		h.Logger.ErrorwCtx(ctx, ErrMissingActorUserID)
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	actorRoles := extractRolesFromContext(ctx)

	span.AddEvent(EventUserIDExtracted,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
			attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
			attribute.StringSlice(commonkeys.Roles, actorRoles),
		),
	)

	cmd := input.DemoteFromAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  actorRoles,
	}

	span.AddEvent(EventAdminServiceDemoteFromAdmin)
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
	)

	// Call service
	updatedUser, err := h.AdminService.DemoteFromAdmin(ctx, cmd)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrDemoteFromAdminFailed, h.Logger)
		return
	}

	h.Logger.InfowCtx(ctx, MsgUserDemotedFromAdmin,
		commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10),
		commonkeys.Username, updatedUser.Username,
		commonkeys.Roles, updatedUser.Roles,
		LogKeyActorUserID, actorUserID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventDemoteFromAdminSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
			attribute.StringSlice(commonkeys.Roles, updatedUser.Roles),
		),
	)

	span.SetStatus(codes.Ok, StatusDemotedFromAdmin)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserRolesResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Roles:     updatedUser.Roles,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUserDemotedFromAdmin)
}

// BlockUser blocks a user by assigning blocked role.
//
// @Summary      Block user
// @Description  Blocks a user by assigning blocked role. Cannot block users with equal or higher privilege.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        user_id  path      integer                      true  "User ID (uint64)"
// @Success      200      {object}  dto.UpdateUserRolesResponse  "User blocked"
// @Failure      400      {string}  string                       "Invalid user_id"
// @Failure      401      {string}  string                       "Unauthorized"
// @Failure      403      {string}  string                       "Forbidden - cannot block user with equal/higher privilege"
// @Failure      404      {string}  string                       "User not found"
// @Failure      500      {string}  string                       "Internal server error"
// @Router       /admin/users/{user_id}/block [put]
//
//nolint:dupl // role management handlers follow similar pattern.
func (h *Handler) BlockUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerName).
		Start(r.Context(), SpanBlockUser)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	// Extract target user ID from path
	userIDParam := chi.URLParam(r, commonkeys.UserID)
	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, ErrMissingUserIDParam)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam, commonkeys.Error, err.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, ErrInvalidUserIDParam)
		span.RecordError(validationErr)
		span.SetStatus(codes.Error, ErrInvalidUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, validationErr, h.Logger)
		return
	}

	// Extract actor user ID and roles from context
	actorUserID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok {
		err := sharederrors.ErrUnauthorized(ErrMissingActorUserID)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingActorUserID)
		h.Logger.ErrorwCtx(ctx, ErrMissingActorUserID)
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	actorRoles := extractRolesFromContext(ctx)

	span.AddEvent(EventUserIDExtracted,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
			attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
			attribute.StringSlice(commonkeys.Roles, actorRoles),
		),
	)

	cmd := input.BlockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  actorRoles,
	}

	span.AddEvent(EventAdminServiceBlockUser)
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
	)

	// Call service
	updatedUser, err := h.AdminService.BlockUser(ctx, cmd)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrBlockUserFailed, h.Logger)
		return
	}

	h.Logger.InfowCtx(ctx, MsgUserBlocked,
		commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10),
		commonkeys.Username, updatedUser.Username,
		commonkeys.Roles, updatedUser.Roles,
		LogKeyActorUserID, actorUserID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventBlockUserSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
			attribute.StringSlice(commonkeys.Roles, updatedUser.Roles),
		),
	)

	span.SetStatus(codes.Ok, StatusUserBlocked)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserRolesResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Roles:     updatedUser.Roles,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUserBlocked)
}

// UnblockUser unblocks a user by restoring default user role.
//
// @Summary      Unblock user
// @Description  Unblocks a user by restoring default user role.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        user_id  path      integer                      true  "User ID (uint64)"
// @Success      200      {object}  dto.UpdateUserRolesResponse  "User unblocked"
// @Failure      400      {string}  string                       "Invalid user_id"
// @Failure      401      {string}  string                       "Unauthorized"
// @Failure      403      {string}  string                       "Forbidden - insufficient privileges"
// @Failure      404      {string}  string                       "User not found"
// @Failure      500      {string}  string                       "Internal server error"
// @Router       /admin/users/{user_id}/unblock [put]
//
//nolint:dupl // role management handlers follow similar pattern.
func (h *Handler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerName).
		Start(r.Context(), SpanUnblockUser)
	defer span.End()

	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	span.AddEvent(EventDecodeRequest,
		trace.WithAttributes(
			attribute.String(tracingkeys.RequestIPKey, ip),
			attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		),
	)

	// Extract target user ID from path
	userIDParam := chi.URLParam(r, commonkeys.UserID)
	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, ErrMissingUserIDParam)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam, commonkeys.Error, err.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, ErrInvalidUserIDParam)
		span.RecordError(validationErr)
		span.SetStatus(codes.Error, ErrInvalidUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		httpresponse.WriteValidationErrorSpan(ctx, w, span, validationErr, h.Logger)
		return
	}

	// Extract actor user ID and roles from context
	actorUserID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok {
		err := sharederrors.ErrUnauthorized(ErrMissingActorUserID)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMissingActorUserID)
		h.Logger.ErrorwCtx(ctx, ErrMissingActorUserID)
		httpresponse.WriteAuthError(w, err, h.Logger)
		return
	}

	actorRoles := extractRolesFromContext(ctx)

	span.AddEvent(EventUserIDExtracted,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
			attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
			attribute.StringSlice(commonkeys.Roles, actorRoles),
		),
	)

	cmd := input.UnblockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  actorRoles,
	}

	span.AddEvent(EventAdminServiceUnblockUser)
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(targetUserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(actorUserID, 10)),
	)

	// Call service
	updatedUser, err := h.AdminService.UnblockUser(ctx, cmd)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, ErrUnblockUserFailed, h.Logger)
		return
	}

	h.Logger.InfowCtx(ctx, MsgUserUnblocked,
		commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10),
		commonkeys.Username, updatedUser.Username,
		commonkeys.Roles, updatedUser.Roles,
		LogKeyActorUserID, actorUserID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	span.AddEvent(EventUnblockUserSuccess,
		trace.WithAttributes(
			attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
			attribute.StringSlice(commonkeys.Roles, updatedUser.Roles),
		),
	)

	span.SetStatus(codes.Ok, StatusUserUnblocked)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserRolesResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Roles:     updatedUser.Roles,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, MsgUserUnblocked)
}

// extractRolesFromContext extracts roles from the context claims.
func extractRolesFromContext(ctx context.Context) []string {
	claimsVal := ctx.Value(ctxkeys.Claims)
	if claimsVal == nil {
		return []string{}
	}

	claims, ok := claimsVal.(map[string]any)
	if !ok {
		return []string{}
	}

	rolesVal, exists := claims[commonkeys.Roles]
	if !exists {
		return []string{}
	}

	// Handle []string
	if roles, ok := rolesVal.([]string); ok {
		return roles
	}

	// Handle []any (interface slice)
	if rolesAny, ok := rolesVal.([]any); ok {
		roles := make([]string, 0, len(rolesAny))
		for _, r := range rolesAny {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
		return roles
	}

	return []string{}
}
