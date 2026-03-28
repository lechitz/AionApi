package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Clear removes chat history cache for a user.
func (s *Store) Clear(ctx context.Context, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameChatHistoryClear, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	s.logger.InfowCtx(ctx, LogClearHistory, commonkeys.UserID, userID)

	key := buildKey(userID)
	if err := s.cache.Del(ctx, key); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrClearHistory)
		s.logger.ErrorwCtx(ctx, ErrClearHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return fmt.Errorf("%s: %w", ErrClearHistory, err)
	}

	span.SetStatus(codes.Ok, LogHistoryCleared)
	s.logger.InfowCtx(ctx, LogHistoryCleared, commonkeys.UserID, userID)

	return nil
}
