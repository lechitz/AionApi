// Package repository provides methods for interacting with the category database.
package repository

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"

	"gorm.io/gorm"
)

// CategoryRepository manages database operations related to category entities.
// It uses gorm.DB for ORM and output.ContextLogger for logging operations.
type CategoryRepository struct {
	db     *gorm.DB
	logger output.ContextLogger
}

// NewCategory creates a new instance of CategoryRepository with a given gorm.DB and contextlogger.
func NewCategory(db *gorm.DB, logger output.ContextLogger) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}
