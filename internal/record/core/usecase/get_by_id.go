package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// GetByID returns a record by id for the user.
func (s *Service) GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error) {
	rec, err := s.RecordRepository.GetByID(ctx, recordID, userID)
	if err != nil {
		return domain.Record{}, err
	}
	return rec, nil
}
