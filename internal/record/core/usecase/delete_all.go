package usecase

import "context"

// DeleteAll performs a soft delete on all records for the given user.
func (s *Service) DeleteAll(ctx context.Context, userID uint64) error {
	return s.RecordRepository.DeleteAllByUser(ctx, userID)
}
