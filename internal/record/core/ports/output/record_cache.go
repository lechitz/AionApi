// Package output defines interfaces for record-related cache operations.
package output

import (
	"context"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
)

// RecordCache defines cache operations for records.
type RecordCache interface {
	SaveRecord(ctx context.Context, record domain.Record, expiration time.Duration) error
	SaveRecordsByDay(ctx context.Context, userID uint64, date time.Time, records []domain.Record, expiration time.Duration) error
	SaveRecordsByCategory(ctx context.Context, categoryID, userID uint64, records []domain.Record, expiration time.Duration) error
	SaveRecordsByTag(ctx context.Context, tagID, userID uint64, records []domain.Record, expiration time.Duration) error
	GetRecord(ctx context.Context, recordID, userID uint64) (domain.Record, error)
	GetRecordsByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error)
	GetRecordsByCategory(ctx context.Context, categoryID, userID uint64) ([]domain.Record, error)
	GetRecordsByTag(ctx context.Context, tagID, userID uint64) ([]domain.Record, error)
	DeleteRecord(ctx context.Context, recordID, userID uint64) error
	DeleteRecordsByDay(ctx context.Context, userID uint64, date time.Time) error
	DeleteRecordsByCategory(ctx context.Context, categoryID, userID uint64) error
	DeleteRecordsByTag(ctx context.Context, tagID, userID uint64) error
}
