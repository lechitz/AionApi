package usecase

import (
	"sync"
	"sync/atomic"

	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/realtime/core/domain"
)

// Service provides fan-out of realtime events to per-user subscribers.
type Service struct {
	logger           logger.ContextLogger
	subscriberBuffer int

	mu          sync.RWMutex
	nextID      uint64
	subscribers map[uint64]map[uint64]chan domain.Event
}

// NewService creates a realtime service with the configured subscriber buffer.
func NewService(log logger.ContextLogger, subscriberBuffer int) *Service {
	if subscriberBuffer <= 0 {
		subscriberBuffer = 1
	}

	return &Service{
		logger:           log,
		subscriberBuffer: subscriberBuffer,
		subscribers:      make(map[uint64]map[uint64]chan domain.Event),
	}
}

func (s *Service) nextSubscriberID() uint64 {
	return atomic.AddUint64(&s.nextID, 1)
}
