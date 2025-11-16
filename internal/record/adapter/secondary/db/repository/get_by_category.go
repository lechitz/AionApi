package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// ListByCategory returns records filtered by category for a given user.
func (r *RecordRepository) ListByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
	var recordsDB []model.Record
	q := r.db.WithContext(ctx).
		Where("category_id = ? AND user_id = ? AND deleted_at IS NULL", categoryID, userID).
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
