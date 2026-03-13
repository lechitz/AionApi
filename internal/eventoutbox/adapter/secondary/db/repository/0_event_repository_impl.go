package repository

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// EventRepository manages DB operations for durable outbox events.
type EventRepository struct {
	db     db.DB
	logger logger.ContextLogger
}

// NewEventRepository creates a new durable outbox event repository.
func NewEventRepository(database db.DB, log logger.ContextLogger) *EventRepository {
	return &EventRepository{
		db:     database,
		logger: log,
	}
}

// WithDB clones the repository with a transaction-bound database handle.
func (r *EventRepository) WithDB(database db.DB) *EventRepository {
	if r == nil {
		return nil
	}
	return &EventRepository{
		db:     database,
		logger: r.logger,
	}
}
