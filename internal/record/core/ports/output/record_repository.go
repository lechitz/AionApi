// Package output defines persistence output ports for the record bounded context.
package output

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// RecordRepository defines persistence operations for records.
type RecordRepository interface {
	Create(ctx context.Context, r domain.Record) (domain.Record, error)
	Update(ctx context.Context, r domain.Record) (domain.Record, error)
	GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error)
	GetByUserCategoryDate(ctx context.Context, userID uint64, categoryID uint64, date time.Time) (domain.Record, error)

	ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	ListByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int) ([]domain.Record, error)
	ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error)
	ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error)
	ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error)

	Delete(ctx context.Context, id uint64, userID uint64) error
	DeleteAllByUser(ctx context.Context, userID uint64) error
}
