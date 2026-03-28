package repository

import (
	"context"

	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
)

// ListByCategory returns records filtered by category for a given user.
// Records are retrieved via subquery filtering tag_ids that belong to the category.
func (r *RecordRepository) ListByCategory(
	ctx context.Context,
	categoryID uint64,
	userID uint64,
	limit int,
	afterEventTime *string,
	afterID *int64,
) ([]domain.Record, error) {
	var recordsDB []model.Record

	// Subquery: get all tag_ids for this category
	q := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Where("tag_id IN (SELECT tag_id FROM aion_api.tags WHERE category_id = ?)", categoryID).
		Order("event_time DESC, id DESC").
		Limit(limit)

	if afterEventTime != nil && afterID != nil {
		q = q.Where("event_time < ? OR (event_time = ? AND id < ?)", *afterEventTime, *afterEventTime, *afterID)
	}

	if err := q.Find(&recordsDB).Error(); err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}
