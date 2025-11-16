package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByDay returns all records for a specific day (based on event_time date) for a given user.
func (r *RecordRepository) ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error) {
	var recordsDB []model.Record
	// Extract date boundaries (00:00:00 to 23:59:59)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND event_time >= ? AND event_time < ? AND deleted_at IS NULL", userID, startOfDay, endOfDay).
		Order("event_time ASC, id ASC").
		Find(&recordsDB).Error; err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
