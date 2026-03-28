// Package input defines use case interfaces for the event outbox context.
package input

import (
	"context"

	"github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
)

// Service defines durable outbox write operations.
type Service interface {
	Enqueue(ctx context.Context, event domain.Event) error
}

// PublisherService defines background publication of pending outbox rows.
type PublisherService interface {
	PublishPending(ctx context.Context, limit int) error
}
