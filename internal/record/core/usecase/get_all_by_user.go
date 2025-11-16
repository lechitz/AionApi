package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByUser returns records for a user with optional cursor parameters.
func (s *Service) ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
	return s.RecordRepository.ListByUser(ctx, userID, limit, afterEventTime, afterID)
}
