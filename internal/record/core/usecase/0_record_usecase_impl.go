// Package usecase implements application business logic for records.
package usecase

import (
	"context"
	"errors"
	"strconv"

	eventoutboxinput "github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/record/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	tagoutput "github.com/lechitz/AionApi/internal/tag/core/ports/output"
)

// Service implements the record use cases.
type Service struct {
	RecordRepository output.RecordRepository
	RecordCache      output.RecordCache
	TagRepository    tagoutput.TagRepository
	OutboxService    eventoutboxinput.Service
	Logger           logger.ContextLogger
}

// NewService is a convention wrapper used by bootstrap to instantiate the record service.
func NewService(recordRepo output.RecordRepository, cache output.RecordCache, tagRepo tagoutput.TagRepository, logger logger.ContextLogger) *Service {
	return &Service{
		RecordRepository: recordRepo,
		RecordCache:      cache,
		TagRepository:    tagRepo,
		Logger:           logger,
	}
}

// WithOutbox attaches an optional outbox service without breaking existing constructor call sites.
func (s *Service) WithOutbox(outboxService eventoutboxinput.Service) *Service {
	s.OutboxService = outboxService
	return s
}

// getUserIDFromContext extracts a numeric user ID from context, supporting common types.
func getUserIDFromContext(ctx context.Context) (uint64, error) {
	v := ctx.Value(ctxkeys.UserID)
	if v == nil {
		return 0, errors.New(UserNotAuthenticated)
	}
	switch id := v.(type) {
	case uint64:
		return id, nil
	case int64:
		if id < 0 {
			return 0, errors.New(UserIDNegative)
		}
		return uint64(id), nil
	case int:
		if id < 0 {
			return 0, errors.New(UserIDNegative)
		}
		return uint64(id), nil
	case string:
		// try parse as uint64
		if u, err := strconv.ParseUint(id, 10, 64); err == nil {
			return u, nil
		}
		return 0, errors.New(UserIDStringNotSupported)
	default:
		return 0, errors.New(InvalidUserIDInContext)
	}
}
