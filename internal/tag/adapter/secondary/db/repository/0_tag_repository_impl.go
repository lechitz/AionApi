// Package repository provides methods for interacting with the handler database.
package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// TagRepository manages database operations related to tag entities.
// Depends on db.DB interface (not *gorm.DB) following Hexagonal Architecture.
type TagRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// New creates a new instance of TagRepository with a given database connection and logger.
func New(database db.DB, logger logger.ContextLogger) *TagRepository {
	return &TagRepository{
		db:     database,
		logger: logger,
	}
}
