package cache

import (
	"context"
	"fmt"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetByKey retrieves a token using a custom cache key.
// This is used for validating tokens during grace period.
func (s *Store) GetByKey(ctx context.Context, key string) (domain.Auth, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTokenGetByKey, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
		attribute.String(AttributeCacheKey, key),
	))
	defer span.End()

	tokenValue, err := s.cache.Get(ctx, key)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Debugw(ErrorToGetTokenFromRedis, AttributeCacheKey, key, commonkeys.Error, err)
		return domain.Auth{}, err
	}

	if tokenValue == "" {
		span.SetStatus(codes.Ok, ErrorTokenNotFoundInGracePeriod)
		return domain.Auth{}, fmt.Errorf(ErrorTokenNotFoundInGracePeriod)
	}

	// Return a generic Auth with the token value
	// The Key and Type are not critical for grace period validation
	tokenDomain := domain.Auth{
		Token: tokenValue,
	}

	span.SetStatus(codes.Ok, TokenRetrievedByCustomKey)
	s.logger.Debugw(TokenRetrievedByCustomKey, AttributeCacheKey, key)
	return tokenDomain, nil
}
