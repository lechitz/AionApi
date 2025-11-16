package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByTag returns records filtered by tag for a given user.
func (r *RecordRepository) ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int) ([]domain.Record, error) {
	var recordsDB []model.Record

	if err := r.db.WithContext(ctx).
		Where("tag_id = ? AND user_id = ? AND deleted_at IS NULL", tagID, userID).
		Order("event_time DESC, id DESC").
		Limit(limit).
		Find(&recordsDB).Error; err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
