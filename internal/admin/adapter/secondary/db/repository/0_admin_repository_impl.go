// Package repository provides methods for interacting with the admin database.
package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// AdminRepository handles admin-related database operations.
// Depends on db.DB interface following Hexagonal Architecture.
type AdminRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// New initializes a new AdminRepository with the provided database connection and logger.
func New(database db.DB, logger logger.ContextLogger) *AdminRepository {
	return &AdminRepository{
		db:     database,
		logger: logger,
	}
}
