package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Get returns the token for a user or empty if not present.
func (s *Store) Get(ctx context.Context, tokenKey uint64) (domain.Auth, error) {
	tr := otel.Tracer(SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, SpanNameTokenGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(tokenKey, 10)),
		attribute.String(commonkeys.Operation, commonkeys.OperationGet),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TokenUserKeyFormat, tokenKey)

	tokenValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetTokenFromRedis, commonkeys.TokenKey, cacheKey, commonkeys.Error, err)
		return domain.Auth{}, err
	}

	if tokenValue == "" {
		span.SetStatus(codes.Ok, sharederrors.ErrTokenNotFound)
		return domain.Auth{}, nil
	}

	tokenDomain := domain.Auth{
		Key:   tokenKey,
		Token: tokenValue,
	}

	span.SetStatus(codes.Ok, TokenRetrievedSuccessfully)
	return tokenDomain, nil
}
