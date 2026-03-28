package repository

import (
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// AuditActionEventRepository manages DB operations for immutable audit action events.
type AuditActionEventRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// NewAuditActionEventRepository creates a new audit action event repository.
func NewAuditActionEventRepository(database db.DB, log logger.ContextLogger) *AuditActionEventRepository {
	return &AuditActionEventRepository{
		db:     database,
		logger: log,
	}
}
