package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByUser returns records for the given user ordered by event_time desc.
func (r *RecordRepository) ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
	var recordsDB []model.Record
	q := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("event_time DESC, id DESC").
		Limit(limit)

	if afterEventTime != nil && afterID != nil {
		q = q.Where("event_time < ? OR (event_time = ? AND id < ?)", *afterEventTime, *afterEventTime, *afterID)
	}

	if err := q.Find(&recordsDB).Error; err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
