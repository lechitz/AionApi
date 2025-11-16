package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListAllBetween returns records with event_time within the specified date range (inclusive).
func (s *Service) ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("startDate must be before or equal to endDate")
	}

	records, err := s.RecordRepository.ListAllBetween(ctx, userID, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	return records, nil
}
