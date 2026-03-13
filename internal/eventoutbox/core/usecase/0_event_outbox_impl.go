package usecase

import (
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Service provides durable enqueue operations for canonical outbox events.
type Service struct {
	repository output.EventRepository
	logger     logger.ContextLogger
	now        func() time.Time
}

// NewService creates a new outbox service implementation.
func NewService(repository output.EventRepository, log logger.ContextLogger) input.Service {
	return &Service{
		repository: repository,
		logger:     log,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}
