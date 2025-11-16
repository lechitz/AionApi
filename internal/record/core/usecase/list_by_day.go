package usecase

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByDay returns all records for a specific day for the authenticated user.
func (s *Service) ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error) {
	records, err := s.RecordRepository.ListByDay(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	return records, nil
}
