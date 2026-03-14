package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/realtime/core/domain"
)

type Service interface {
	Publish(ctx context.Context, event domain.Event)
	Subscribe(ctx context.Context, userID uint64) (<-chan domain.Event, func())
}
