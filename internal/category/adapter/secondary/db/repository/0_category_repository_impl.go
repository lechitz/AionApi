// Package repository provides methods for interacting with the handler database.
package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// CategoryRepository manages database operations related to category entities.
// Depends on db.DB interface (not *gorm.DB) following Hexagonal Architecture.
type CategoryRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// New creates a new instance of CategoryRepository with a given database connection and logger.
func New(database db.DB, logger logger.ContextLogger) *CategoryRepository {
	return &CategoryRepository{
		db:     database,
		logger: logger,
	}
}
