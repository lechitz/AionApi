// Package repository contains DB adapters for the record bounded context.
package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"gorm.io/gorm"
)

// RecordRepository implements persistence using GORM and Postgres.
type RecordRepository struct {
	db     *gorm.DB
	logger logger.ContextLogger
}

// New creates a new RecordRepository.
func New(db *gorm.DB, logger logger.ContextLogger) *RecordRepository {
	return &RecordRepository{db: db, logger: logger}
}
