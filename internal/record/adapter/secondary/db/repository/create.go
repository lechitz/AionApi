package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// Create inserts a record and returns the created entity with ID populated.
func (r *RecordRepository) Create(ctx context.Context, rec domain.Record) (domain.Record, error) {
	recordDB := mapper.RecordToDB(rec)

	if err := r.db.WithContext(ctx).Create(&recordDB).Error; err != nil {
		return domain.Record{}, err
	}

	return mapper.RecordFromDB(recordDB), nil
}
