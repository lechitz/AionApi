package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// RecordProjectionRetriever defines derived projection read operations.
type RecordProjectionRetriever interface {
	GetProjectedByID(ctx context.Context, recordID uint64, userID uint64) (domain.RecordProjection, error)
	ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]domain.RecordProjection, error)
}
