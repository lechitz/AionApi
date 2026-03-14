// Package repository provides methods for interacting with the admin database.
package repository

import (
	"context"
	"strconv"
	"time"

	adminmapper "github.com/lechitz/AionApi/internal/admin/adapter/secondary/db/mapper"
	adminmodel "github.com/lechitz/AionApi/internal/admin/adapter/secondary/db/model"
	admindomain "github.com/lechitz/AionApi/internal/admin/core/domain"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	usermapper "github.com/lechitz/AionApi/internal/user/adapter/secondary/db/mapper"
	usermodel "github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetByID retrieves a user by their ID including their roles from user_roles table.
// Returns AdminUser (admin bounded context domain) instead of user.User.
func (r *AdminRepository) GetByID(ctx context.Context, userID uint64) (admindomain.AdminUser, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByID, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGetByID),
	))
	defer span.End()

	r.logger.InfowCtx(ctx, LogFetchingUser, commonkeys.UserID, userID)

	var userDB usermodel.UserDB
	if err := r.db.WithContext(ctx).
		Where(commonkeys.UserID+" = ?", userID).
		First(&userDB).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, LogFailedGetUser, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return admindomain.AdminUser{}, err
	}

	// Get roles from user_roles junction table
	roles, err := r.getUserRoles(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, LogFailedGetUser, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return admindomain.AdminUser{}, err
	}

	// Convert to user domain first (repository layer)
	userDomain := usermapper.UserFromDB(userDB)

	// Convert user.User to admin.AdminUser (anti-corruption layer)
	adminUser := adminmapper.AdminUserFromUser(userDomain)
	adminUser.Roles = roles

	span.SetStatus(codes.Ok, StatusUserFetched)
	r.logger.InfowCtx(ctx, LogUserFetched, commonkeys.UserID, userID)

	return adminUser, nil
}

// getUserRoles fetches the role names for a user from the junction table.
func (r *AdminRepository) getUserRoles(ctx context.Context, userID uint64) ([]string, error) {
	var roles []adminmodel.RoleDB

	// Query joining user_roles with roles table using Raw SQL
	query := `
		SELECT r.role_id, r.name, r.description, r.is_active, r.created_at, r.updated_at
		FROM aion_api.user_roles ur
		JOIN aion_api.roles r ON r.role_id = ur.role_id
		WHERE ur.user_id = ? AND r.is_active = true
	`

	if err := r.db.WithContext(ctx).
		Raw(query, userID).
		Scan(&roles).Error(); err != nil {
		return nil, err
	}

	// Extract role names
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	// Default to "user" if no roles found
	if len(roleNames) == 0 {
		return []string{"user"}, nil
	}

	return roleNames, nil
}

// UpdateRoles updates the roles for a user using the user_roles junction table.
// Returns AdminUser (admin bounded context domain).
func (r *AdminRepository) UpdateRoles(ctx context.Context, userID uint64, roles []string) (admindomain.AdminUser, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdateRoles, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationUpdateRoles),
		attribute.StringSlice(commonkeys.Roles, roles),
	))
	defer span.End()

	r.logger.InfowCtx(ctx, LogUpdatingRoles,
		commonkeys.UserID, userID,
		commonkeys.Roles, roles,
	)

	// Use closure-based transaction (auto-commit on success, auto-rollback on error)
	err := r.db.WithContext(ctx).Transaction(func(tx db.DB) error {
		// Delete existing roles for user
		if err := tx.Where("user_id = ?", userID).
			Delete(&adminmodel.UserRoleDB{}).Error(); err != nil {
			span.RecordError(err)
			r.logger.ErrorwCtx(ctx, LogFailedUpdateRoles,
				commonkeys.UserID, userID,
				commonkeys.Error, err.Error(),
			)
			return err
		}

		// Get role IDs for the new roles
		var roleDBs []adminmodel.RoleDB
		if err := tx.Where("name IN ?", roles).
			Where("is_active = ?", true).
			Find(&roleDBs).Error(); err != nil {
			span.RecordError(err)
			r.logger.ErrorwCtx(ctx, LogFailedUpdateRoles,
				commonkeys.UserID, userID,
				commonkeys.Error, err.Error(),
			)
			return err
		}

		// Insert new role assignments
		now := time.Now().UTC()
		for _, roleDB := range roleDBs {
			userRole := adminmodel.UserRoleDB{
				UserID:     userID,
				RoleID:     roleDB.ID,
				AssignedAt: now,
			}
			if err := tx.Create(&userRole).Error(); err != nil {
				span.RecordError(err)
				r.logger.ErrorwCtx(ctx, LogFailedUpdateRoles,
					commonkeys.UserID, userID,
					commonkeys.Error, err.Error(),
				)
				return err
			}
		}

		// Update user's updated_at timestamp
		if err := tx.Model(&usermodel.UserDB{}).
			Where(commonkeys.UserID+" = ?", userID).
			Update(commonkeys.UserUpdatedAt, now).Error(); err != nil {
			span.RecordError(err)
			r.logger.ErrorwCtx(ctx, LogFailedUpdateRoles,
				commonkeys.UserID, userID,
				commonkeys.Error, err.Error(),
			)
			return err
		}

		return nil // Success - auto commit
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return admindomain.AdminUser{}, err
	}

	span.SetStatus(codes.Ok, StatusRolesUpdated)
	r.logger.InfowCtx(ctx, LogRolesUpdated, commonkeys.UserID, userID, commonkeys.Roles, roles)

	return r.GetByID(ctx, userID)
}
