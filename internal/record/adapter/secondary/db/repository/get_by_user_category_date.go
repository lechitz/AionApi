package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// GetByUserCategoryDate retrieves a record by user+category+date
func (r *RecordRepository) GetByUserCategoryDate(ctx context.Context, userID uint64, categoryID uint64, date time.Time) (domain.Record, error) {
	var recordDB model.Record

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ? AND event_time = ? AND deleted_at IS NULL", userID, categoryID, date).
		First(&recordDB).Error; err != nil {
		return domain.Record{}, err
	}

	return mapper.RecordFromDB(recordDB), nil
}
