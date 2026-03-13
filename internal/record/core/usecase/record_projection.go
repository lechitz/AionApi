package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// GetProjectedByID returns one derived record projection for the user.
func (s *Service) GetProjectedByID(ctx context.Context, recordID uint64, userID uint64) (domain.RecordProjection, error) {
	if userID == 0 {
		return domain.RecordProjection{}, ErrUserIDIsRequired
	}
	if recordID == 0 {
		return domain.RecordProjection{}, ErrRecordIDIsRequired
	}
	if s.RecordProjectionRepository == nil {
		return domain.RecordProjection{}, ErrProjectionRepositoryUnavailable
	}

	return s.RecordProjectionRepository.GetProjectedByID(ctx, userID, recordID)
}

// ListProjectedLatest returns the latest derived record projections for the user.
func (s *Service) ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]domain.RecordProjection, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	if s.RecordProjectionRepository == nil {
		return nil, ErrProjectionRepositoryUnavailable
	}

	if limit <= 0 {
		limit = 10
	}
	return s.RecordProjectionRepository.ListProjectedLatest(ctx, userID, limit)
}
