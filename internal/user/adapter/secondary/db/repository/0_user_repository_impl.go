// Package repository provides methods for interacting with the user database.
package repository

import (
	"strings"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// UserRepository handles interactions with the user database.
// Depends on db.DB interface (not *gorm.DB) following Hexagonal Architecture.
type UserRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// New initializes a new UserRepository with the provided database connection and logger.
func New(database db.DB, logger logger.ContextLogger) *UserRepository {
	return &UserRepository{
		db:     database,
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
	case strings.Contains(msg, PgConstraintUsersUsernameKey):
		return commonkeys.Username, true
	case strings.Contains(msg, PgConstraintUsersEmailKey):
		return commonkeys.Email, true
	}
	return "", false
}
