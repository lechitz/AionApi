// Package mapper contains functions to convert between domain and database models.
package mapper

import (
	"github.com/lechitz/aion-api/internal/admin/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/admin/core/domain"
)

// RoleFromDB converts a model.RoleDB object into a domain.Role object.
func RoleFromDB(role model.RoleDB) domain.Role {
	var description string
	if role.Description != nil {
		description = *role.Description
	}

	return domain.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// RoleToDB converts a domain.Role object into a model.RoleDB object.
func RoleToDB(role domain.Role) model.RoleDB {
	var description *string
	if role.Description != "" {
		description = &role.Description
	}

	return model.RoleDB{
		ID:          role.ID,
		Name:        role.Name,
		Description: description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// UserRoleFromDB converts a model.UserRoleDB object into a domain.UserRole object.
func UserRoleFromDB(ur model.UserRoleDB) domain.UserRole {
	return domain.UserRole{
		ID:         ur.ID,
		UserID:     ur.UserID,
		RoleID:     ur.RoleID,
		AssignedBy: ur.AssignedBy,
		AssignedAt: ur.AssignedAt,
	}
}

// UserRoleToDB converts a domain.UserRole object into a model.UserRoleDB object.
func UserRoleToDB(ur domain.UserRole) model.UserRoleDB {
	return model.UserRoleDB{
		ID:         ur.ID,
		UserID:     ur.UserID,
		RoleID:     ur.RoleID,
		AssignedBy: ur.AssignedBy,
		AssignedAt: ur.AssignedAt,
	}
}
