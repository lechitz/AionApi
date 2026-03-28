// Package repository contains DB adapters for the record bounded context.
package repository

import (
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
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

// WithDB clones the repository with a transaction-bound database handle.
func (r *RecordRepository) WithDB(database db.DB) *RecordRepository {
	if r == nil {
		return nil
	}
	return &RecordRepository{
		db:     database,
		logger: r.logger,
	}
}
