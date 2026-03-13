// Package output defines output ports for the event outbox context.
package output

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
)

// EventRepository persists canonical outbox events.
type EventRepository interface {
	Save(ctx context.Context, event domain.Event) error
	ListPending(ctx context.Context, limit int) ([]domain.Event, error)
	MarkPublished(ctx context.Context, eventID string, publishedAt time.Time) error
	Reschedule(ctx context.Context, eventID string, nextAvailableAt time.Time, lastError string) error
}
