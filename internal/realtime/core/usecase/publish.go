package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/realtime/core/domain"
)

func (s *Service) Publish(ctx context.Context, event domain.Event) {
	s.mu.RLock()
	subscribers := s.subscribers[event.UserID]
	if len(subscribers) == 0 {
		s.mu.RUnlock()
		return
	}

	channels := make([]chan domain.Event, 0, len(subscribers))
	for _, ch := range subscribers {
		channels = append(channels, ch)
	}
	s.mu.RUnlock()

	for _, ch := range channels {
		select {
		case ch <- event:
		default:
			s.logger.WarnwCtx(ctx, logRealtimeEventDropped,
				"user_id", event.UserID,
				"record_id", event.RecordID,
				"event_type", event.Type,
			)
		}
	}
}
