package usecase

import (
	"time"

	"github.com/lechitz/aion-api/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/aion-api/internal/eventoutbox/core/ports/output"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// PublisherService publishes pending outbox rows to the event backbone.
type PublisherService struct {
	repository output.EventRepository
	publisher  output.EventPublisher
	logger     logger.ContextLogger
	now        func() time.Time
	backoff    time.Duration
}

// NewPublisherService creates a new background outbox publisher service.
func NewPublisherService(
	repository output.EventRepository,
	publisher output.EventPublisher,
	log logger.ContextLogger,
) input.PublisherService {
	return &PublisherService{
		repository: repository,
		publisher:  publisher,
		logger:     log,
		now: func() time.Time {
			return time.Now().UTC()
		},
		backoff: DefaultPublishBackoff,
	}
}
