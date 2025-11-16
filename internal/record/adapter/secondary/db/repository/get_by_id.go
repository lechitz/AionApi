package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// GetByID retrieves a record by ID ensuring it belongs to the provided user.
func (r *RecordRepository) GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error) {
	var recordDB model.Record

	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", recordID, userID).
		First(&recordDB).Error; err != nil {
		return domain.Record{}, err
	}

	return mapper.RecordFromDB(recordDB), nil
}
