package repository

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
)

// Delete performs a soft delete for the given record belonging to the user.
func (r *RecordRepository) Delete(ctx context.Context, recordID uint64, userID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.Record{}).
		Where("id = ? AND user_id = ?", recordID, userID).
		Update("deleted_at", time.Now()).Error
}
