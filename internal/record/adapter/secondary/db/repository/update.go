package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// Update updates a record and returns the updated entity.
func (r *RecordRepository) Update(ctx context.Context, rec domain.Record) (domain.Record, error) {
	recDB := mapper.RecordToDB(rec)
	if err := r.db.WithContext(ctx).Model(&model.Record{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", rec.ID, rec.UserID).
		Updates(recDB).Error; err != nil {
		return domain.Record{}, err
	}

	var out model.Record
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", rec.ID, rec.UserID).
		First(&out).Error; err != nil {
		return domain.Record{}, err
	}
	return mapper.RecordFromDB(out), nil
}
