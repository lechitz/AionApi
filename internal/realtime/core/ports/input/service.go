// Package input exposes realtime use case contracts.
package input

import (
	"context"

	"github.com/lechitz/aion-api/internal/realtime/core/domain"
)

// Service defines the realtime publish and subscribe operations.
type Service interface {
	Publish(ctx context.Context, event domain.Event)
	Subscribe(ctx context.Context, userID uint64) (<-chan domain.Event, func())
}
