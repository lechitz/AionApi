// Package repository contains DB adapters for the record bounded context.
package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// RecordRepository implements persistence using database operations.
// Depends on db.DB interface (not *gorm.DB) following Hexagonal Architecture.
type RecordRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// New creates a new RecordRepository.
func New(database db.DB, logger logger.ContextLogger) *RecordRepository {
	return &RecordRepository{db: database, logger: logger}
}
