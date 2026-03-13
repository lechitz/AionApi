package usecase

import (
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/ports/output"
	dbport "github.com/lechitz/AionApi/internal/platform/ports/output/db"
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

type eventRepositoryWithDB interface {
	WithDB(database dbport.DB) output.EventRepository
}

// WithDB clones the outbox service with a transaction-bound repository when supported.
func (s *Service) WithDB(database dbport.DB) input.Service {
	if s == nil {
		return nil
	}

	if repository, ok := s.repository.(eventRepositoryWithDB); ok {
		return &Service{
			repository: repository.WithDB(database),
			logger:     s.logger,
			now:        s.now,
		}
	}

	return s
}
