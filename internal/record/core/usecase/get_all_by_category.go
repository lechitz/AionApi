package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListRecordsByCategory lists records filtered by category.
func (s *Service) ListRecordsByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int, after *string) ([]domain.Record, error) {
	_ = after
	return s.RecordRepository.ListByCategory(ctx, categoryID, userID, limit, nil, nil)
}
