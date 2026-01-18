package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListLatest returns the N most recent records for a user, ordered by event_time DESC.
func (r *RecordRepository) ListLatest(ctx context.Context, userID uint64, limit int) ([]domain.Record, error) {
	var recordsDB []model.Record

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("event_time DESC, id DESC").
		Limit(limit).
		Find(&recordsDB).Error(); err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
