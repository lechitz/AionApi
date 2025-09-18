// Package repository provides methods for interacting with the handler database.
package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"gorm.io/gorm"
)

// CategoryRepository manages database operations related to category entities.
// It uses gorm.DB for ORM and output.ContextLogger for logging operations.
type CategoryRepository struct {
	db     *gorm.DB
	logger logger.ContextLogger
}

// New creates a new instance of CategoryRepository with a given gorm.DB and contextlogger.
func New(db *gorm.DB, logger logger.ContextLogger) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}
