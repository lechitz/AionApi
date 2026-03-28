// Package usecase implements admin role management operations.
package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/admin/core/domain"
	"github.com/lechitz/aion-api/internal/admin/core/ports/input"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// PromoteToAdmin promotes a user to admin role.
// Only users with higher hierarchy (owner, admin) can promote to admin.
func (s *Service) PromoteToAdmin(ctx context.Context, cmd input.PromoteToAdminCommand) (domain.AdminUser, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanPromoteToAdmin)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(cmd.ActorUserID, 10)),
	)

	// Check authorization: actor must have admin or owner role
	actorHighestRole := domain.GetHighestRole(cmd.ActorRoles)
	if !domain.CanManageRole(actorHighestRole, domain.RoleAdmin) {
		err := ErrUnauthorizedPromoteToAdmin
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusUnauthorized)
		s.logger.WarnwCtx(ctx, WarnUnauthorizedPromoteAttempt,
			LogKeyActorUserID, cmd.ActorUserID,
			LogKeyActorRole, actorHighestRole,
			LogKeyTargetUserID, cmd.UserID,
		)
		return domain.AdminUser{}, err
	}

	// Get current user
	currentUser, err := s.adminRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusUserNotFound)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrGetUser, err)
	}

	// Add admin role if not already present
	newRoles := currentUser.Roles
	if !domain.HasRole(newRoles, domain.RoleAdmin) {
		newRoles = append(newRoles, domain.RoleAdmin)
	}

	// Ensure user role is also present
	if !domain.HasRole(newRoles, domain.RoleUser) {
		newRoles = append(newRoles, domain.RoleUser)
	}

	// Update roles
	updatedUser, err := s.adminRepository.UpdateRoles(ctx, cmd.UserID, newRoles)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusUpdateFailed)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrPromoteToAdminFailed, err)
	}

	span.SetStatus(codes.Ok, StatusPromotedToAdmin)
	s.logger.InfowCtx(ctx, InfoUserPromotedToAdmin,
		commonkeys.UserID, cmd.UserID,
		LogKeyByUserID, cmd.ActorUserID,
	)

	return updatedUser, nil
}

// DemoteFromAdmin removes admin role from a user.
func (s *Service) DemoteFromAdmin(ctx context.Context, cmd input.DemoteFromAdminCommand) (domain.AdminUser, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanDemoteFromAdmin)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(cmd.ActorUserID, 10)),
	)

	// Check authorization
	actorHighestRole := domain.GetHighestRole(cmd.ActorRoles)
	if !domain.CanManageRole(actorHighestRole, domain.RoleAdmin) {
		err := ErrUnauthorizedDemoteFromAdmin
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusUnauthorized)
		return domain.AdminUser{}, err
	}

	// Get current user
	currentUser, err := s.adminRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		span.RecordError(err)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrGetUser, err)
	}

	// Remove admin role, keep other roles
	newRoles := []string{}
	for _, role := range currentUser.Roles {
		if role != domain.RoleAdmin && role != domain.RoleOwner {
			newRoles = append(newRoles, role)
		}
	}

	// Ensure at least user role
	if len(newRoles) == 0 {
		newRoles = []string{domain.RoleUser}
	}

	updatedUser, err := s.adminRepository.UpdateRoles(ctx, cmd.UserID, newRoles)
	if err != nil {
		span.RecordError(err)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrDemoteFromAdminFailed, err)
	}

	span.SetStatus(codes.Ok, StatusDemotedFromAdmin)
	s.logger.InfowCtx(ctx, InfoUserDemotedFromAdmin,
		commonkeys.UserID, cmd.UserID,
		LogKeyByUserID, cmd.ActorUserID,
	)

	return updatedUser, nil
}

// BlockUser blocks a user by assigning blocked role.
func (s *Service) BlockUser(ctx context.Context, cmd input.BlockUserCommand) (domain.AdminUser, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanBlockUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(cmd.ActorUserID, 10)),
	)

	// Get target user
	targetUser, err := s.adminRepository.GetByID(ctx, cmd.UserID)
	if err != nil {
		span.RecordError(err)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrGetUser, err)
	}

	// Check authorization: cannot block users with higher or equal privilege
	actorHighestRole := domain.GetHighestRole(cmd.ActorRoles)
	targetHighestRole := domain.GetHighestRole(targetUser.Roles)

	if !domain.CanManageRole(actorHighestRole, targetHighestRole) {
		err := ErrUnauthorizedBlockUser
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusUnauthorized)
		s.logger.WarnwCtx(ctx, WarnUnauthorizedBlockAttempt,
			LogKeyActorRole, actorHighestRole,
			LogKeyTargetRole, targetHighestRole,
		)
		return domain.AdminUser{}, err
	}

	// Set blocked role (replaces all other roles)
	newRoles := []string{domain.RoleBlocked}

	updatedUser, err := s.adminRepository.UpdateRoles(ctx, cmd.UserID, newRoles)
	if err != nil {
		span.RecordError(err)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrBlockUserFailed, err)
	}

	span.SetStatus(codes.Ok, StatusUserBlocked)
	s.logger.InfowCtx(ctx, InfoUserBlocked,
		commonkeys.UserID, cmd.UserID,
		LogKeyByUserID, cmd.ActorUserID,
	)

	return updatedUser, nil
}

// UnblockUser unblocks a user by restoring default user role.
func (s *Service) UnblockUser(ctx context.Context, cmd input.UnblockUserCommand) (domain.AdminUser, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUnblockUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
		attribute.String(LogKeyActorUserID, strconv.FormatUint(cmd.ActorUserID, 10)),
	)

	// Check authorization
	actorHighestRole := domain.GetHighestRole(cmd.ActorRoles)
	if !domain.CanManageRole(actorHighestRole, domain.RoleBlocked) {
		err := ErrUnauthorizedUnblockUser
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusUnauthorized)
		return domain.AdminUser{}, err
	}

	// Restore default user role
	newRoles := []string{domain.RoleUser}

	updatedUser, err := s.adminRepository.UpdateRoles(ctx, cmd.UserID, newRoles)
	if err != nil {
		span.RecordError(err)
		return domain.AdminUser{}, fmt.Errorf("%w: %w", ErrUnblockUserFailed, err)
	}

	span.SetStatus(codes.Ok, StatusUserUnblocked)
	s.logger.InfowCtx(ctx, InfoUserUnblocked,
		commonkeys.UserID, cmd.UserID,
		LogKeyByUserID, cmd.ActorUserID,
	)

	return updatedUser, nil
}
