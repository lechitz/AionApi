package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// RecordProjectionRepository exposes derived read-model queries owned by the canonical API.
type RecordProjectionRepository interface {
	GetProjectedByID(ctx context.Context, userID uint64, recordID uint64) (domain.RecordProjection, error)
	ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]domain.RecordProjection, error)
	ListProjectedPage(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.RecordProjection, error)
}
