package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListAllUntil returns records with event_time up to (and including) the given timestamp.
func (r *RecordRepository) ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error) {
	var recordsDB []model.Record

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND event_time <= ? AND deleted_at IS NULL", userID, until).
		Order("event_time DESC, id DESC").
		Limit(limit).
		Find(&recordsDB).Error; err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
