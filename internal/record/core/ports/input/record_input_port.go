package input

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// RecordCreator interface for creating a new record.
type RecordCreator interface {
	Create(ctx context.Context, cmd CreateRecordCommand) (domain.Record, error)
}

// RecordRetriever defines methods for retrieving record details.
type RecordRetriever interface {
	GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error)
	ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int) ([]domain.Record, error)
	ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error)
	ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error)
	ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error)
}

// RecordUpdater defines update operations for a record.
type RecordUpdater interface {
	Update(ctx context.Context, recordID uint64, userID uint64, cmd UpdateRecordCommand) (domain.Record, error)
}

// RecordDeleter defines deletion operations for records.
type RecordDeleter interface {
	Delete(ctx context.Context, recordID uint64, userID uint64) error
	DeleteAll(ctx context.Context, userID uint64) error
}

// RecordService defines the input port used by controllers/handlers to interact with record use cases.
type RecordService interface {
	RecordCreator
	RecordRetriever
	RecordUpdater
	RecordDeleter
}
