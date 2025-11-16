package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByTag returns records filtered by tag for the authenticated user.
func (s *Service) ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int) ([]domain.Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	records, err := s.RecordRepository.ListByTag(ctx, tagID, userID, limit)
	if err != nil {
		return nil, err
	}
	return records, nil
}
