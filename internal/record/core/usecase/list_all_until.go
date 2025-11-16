package usecase

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListAllUntil returns records with event_time up to (and including) the given timestamp.
func (s *Service) ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	records, err := s.RecordRepository.ListAllUntil(ctx, userID, until, limit)
	if err != nil {
		return nil, err
	}
	return records, nil
}
