// Package output defines output ports for the event outbox context.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
)

// EventRepository persists canonical outbox events.
type EventRepository interface {
	Save(ctx context.Context, event domain.Event) error
}
