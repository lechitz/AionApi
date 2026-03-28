package cache

import (
	"context"
	"time"

	"github.com/lechitz/aion-api/internal/auth/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SaveWithKey persists a token using a custom cache key with TTL.
// This is used for grace period tokens to allow custom key formats.
func (s *Store) SaveWithKey(ctx context.Context, key string, token domain.Auth, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTokenSaveWithKey, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
		attribute.String(AttributeCacheKey, key),
		attribute.String(AttributeTTL, expiration.String()),
	))
	defer span.End()

	// if no expiration passed, fall back to default
	if expiration <= 0 {
		expiration = TokenExpirationDefault
	}

	span.AddEvent(EventSaveTokenToCacheByKey)
	if err := s.cache.Set(ctx, key, token.Token, expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveTokenToRedis, AttributeCacheKey, key, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TokenSavedWithCustomKey)
	s.logger.Debugw(TokenSavedWithCustomKey, AttributeCacheKey, key, AttributeTTL, expiration.String())
	return nil
}
