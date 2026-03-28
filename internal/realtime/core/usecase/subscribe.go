package usecase

import (
	"context"
	"sync"

	"github.com/lechitz/aion-api/internal/realtime/core/domain"
)

// Subscribe registers a per-user subscriber stream and returns a cleanup function.
func (s *Service) Subscribe(ctx context.Context, userID uint64) (<-chan domain.Event, func()) {
	subscriberID := s.nextSubscriberID()
	ch := make(chan domain.Event, s.subscriberBuffer)

	s.mu.Lock()
	if _, ok := s.subscribers[userID]; !ok {
		s.subscribers[userID] = make(map[uint64]chan domain.Event)
	}
	s.subscribers[userID][subscriberID] = ch
	s.mu.Unlock()

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			subscribers := s.subscribers[userID]
			if subscribers == nil {
				return
			}
			if existing, ok := subscribers[subscriberID]; ok {
				delete(subscribers, subscriberID)
				close(existing)
			}
			if len(subscribers) == 0 {
				delete(s.subscribers, userID)
			}
		})
	}

	go func() {
		<-ctx.Done()
		cleanup()
	}()

	return ch, cleanup
}
