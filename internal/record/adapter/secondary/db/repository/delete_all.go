package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
)

// DeleteAllByUser performs a bulk soft delete for all records of a user.
func (r *RecordRepository) DeleteAllByUser(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.Record{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Update("deleted_at", time.Now()).Error
}
