// Package usecase (admin) provides operations for managing admin functions.
package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/admin/core/domain"
	"github.com/lechitz/AionApi/internal/admin/core/ports/input"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// UpdateUserRoles updates the roles for a given user. Only admins can call this.
func (s *Service) UpdateUserRoles(ctx context.Context, cmd input.UpdateUserRolesCommand) (domain.AdminUser, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUpdateUserRoles)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateUserRoles),
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
		attribute.StringSlice(commonkeys.Roles, cmd.Roles),
	)

	s.logger.InfowCtx(ctx, InfoUpdatingUserRoles,
		commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
		commonkeys.Roles, cmd.Roles,
	)

	if !cmd.HasUpdates() {
		span.RecordError(ErrNoRolesToUpdate)
		span.SetStatus(codes.Error, ErrorNoRolesToUpdate)
		s.logger.ErrorwCtx(ctx, ErrorNoRolesToUpdate, commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10))
		return domain.AdminUser{}, ErrNoRolesToUpdate
	}

	span.AddEvent(SpanEventValidateRoles)
	for _, role := range cmd.Roles {
		if !domain.IsValidRole(role) {
			span.RecordError(ErrInvalidRole)
			span.SetStatus(codes.Error, ErrorInvalidRole)
			s.logger.WarnwCtx(ctx, WarnInvalidRoleProvided,
				commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
				commonkeys.Role, role,
			)
			return domain.AdminUser{}, fmt.Errorf("%w: %s", ErrInvalidRole, role)
		}
	}

	span.AddEvent(SpanEventGetUser)
	currentUser, err := s.adminRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToGetUser)
		s.logger.ErrorwCtx(ctx, ErrorToGetUser,
			commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
			commonkeys.Error, err.Error(),
		)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrGetUser, err)
	}

	if containsRole(cmd.Roles, domain.RoleBlocked) && containsRole(currentUser.Roles, domain.RoleAdmin) {
		span.RecordError(ErrCannotBlockAdmin)
		span.SetStatus(codes.Error, ErrorCannotBlockAdmin)
		s.logger.ErrorwCtx(ctx, ErrorCannotBlockAdmin,
			commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
		)
		return domain.AdminUser{}, ErrCannotBlockAdmin
	}

	span.AddEvent(SpanEventUpdateRoles)
	updatedUser, err := s.adminRepository.UpdateRoles(ctx, cmd.UserID, cmd.Roles)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToUpdateUserRoles)
		s.logger.ErrorwCtx(ctx, ErrorToUpdateUserRoles,
			commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10),
			commonkeys.Error, err.Error(),
		)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrUpdateUserRoles, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
		attribute.StringSlice(commonkeys.Roles, updatedUser.Roles),
	)
	span.SetStatus(codes.Ok, StatusRolesUpdated)

	s.logger.InfowCtx(ctx, InfoUserRolesUpdated,
		commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10),
		commonkeys.Roles, updatedUser.Roles,
	)

	return updatedUser, nil
}

// containsRole checks if a role exists in the roles slice.
func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
