package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListAllBetween returns records with event_time within the specified date range (inclusive).
func (r *RecordRepository) ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error) {
	var recordsDB []model.Record

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND event_time >= ? AND event_time <= ? AND deleted_at IS NULL", userID, startDate, endDate).
		Order("event_time DESC, id DESC").
		Limit(limit).
		Find(&recordsDB).Error; err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
