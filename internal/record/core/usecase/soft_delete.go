package usecase

import "context"

// Delete performs a soft delete for the user's record.
func (s *Service) Delete(ctx context.Context, id uint64, userID uint64) error {
	return s.RecordRepository.Delete(ctx, id, userID)
}
