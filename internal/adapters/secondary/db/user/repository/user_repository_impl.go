// Package repository provides methods for interacting with the user database.
package repository

import (
	"strings"

	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"gorm.io/gorm"
)

// UserRepository handles interactions with the user database.
type UserRepository struct {
	db     *gorm.DB
	logger output.ContextLogger
}

// NewUser initializes a new UserRepository with the provided database connection and logger.
func NewUser(db *gorm.DB, logger output.ContextLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// isUniqueViolation detects a Postgres unique constraint violation and returns the affected field.
func isUniqueViolation(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, constants.PgConstraintUsersUsernameKey):
		return commonkeys.Username, true
	case strings.Contains(msg, constants.PgConstraintUsersEmailKey):
		return commonkeys.Email, true
	}
	return "", false
}
