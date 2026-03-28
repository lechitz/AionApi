// Package output defines output ports for the event outbox context.
package output

import (
	"context"

	"github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
)

// EventPublisher publishes canonical outbox events to the external event backbone.
type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
