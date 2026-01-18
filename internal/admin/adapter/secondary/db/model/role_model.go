// Package model contains the database models for admin entities.
package model

import "time"

const (
	// TableRoles is the name of the database table for roles.
	TableRoles = "aion_api.roles"

	// TableUserRoles is the name of the database table for user_roles.
	TableUserRoles = "aion_api.user_roles"
)

// RoleDB represents the database model for storing role information.
type RoleDB struct {
	ID          uint64    `gorm:"primaryKey;column:role_id"`
	Name        string    `gorm:"column:name;uniqueIndex"`
	Description *string   `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

// TableName specifies the custom database table name for the RoleDB model.
func (RoleDB) TableName() string {
	return TableRoles
}

// UserRoleDB represents the database model for user-role assignments.
type UserRoleDB struct {
	ID         uint64    `gorm:"primaryKey;column:user_role_id"`
	UserID     uint64    `gorm:"column:user_id;index"`
	RoleID     uint64    `gorm:"column:role_id;index"`
	AssignedBy *uint64   `gorm:"column:assigned_by"`
	AssignedAt time.Time `gorm:"column:assigned_at"`
}

// TableName specifies the custom database table name for the UserRoleDB model.
func (UserRoleDB) TableName() string {
	return TableUserRoles
}
