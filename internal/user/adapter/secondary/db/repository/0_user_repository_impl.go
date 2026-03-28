// Package repository provides methods for interacting with the user database.
package repository

import (
	"strings"

	adminoutput "github.com/lechitz/aion-api/internal/admin/core/ports/output"
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
)

// UserRepository handles interactions with the user database.
// Depends on db.DB interface (not *gorm.DB) following Hexagonal Architecture.
// Delegates role assignment to admin context via RoleAssigner interface.
type UserRepository struct {
	db           db.DB
	logger       logger.ContextLogger
	roleAssigner adminoutput.RoleAssigner // Dependency injection from /admin context
}

// New initializes a new UserRepository with the provided dependencies.
func New(database db.DB, logger logger.ContextLogger, roleAssigner adminoutput.RoleAssigner) *UserRepository {
	return &UserRepository{
		db:           database,
		logger:       logger,
		roleAssigner: roleAssigner,
	}
}

// isUniqueViolation detects a Postgres unique constraint violation and returns the affected field.
func isUniqueViolation(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, PgConstraintUsersUsernameKey):
		return commonkeys.Username, true
	case strings.Contains(msg, PgConstraintUsersEmailKey):
		return commonkeys.Email, true
	}
	return "", false
}
