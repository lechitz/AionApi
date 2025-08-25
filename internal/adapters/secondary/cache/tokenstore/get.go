package tokenstore

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/tokenstore/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Get returns the token for a user or empty if not present.
func (s *Store) Get(ctx context.Context, tokenKey uint64) (domain.Token, error) {
	tr := otel.Tracer(constants.SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, constants.SpanNameTokenGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(tokenKey, 10)),
		attribute.String(commonkeys.Operation, commonkeys.OperationGet),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(constants.TokenUserKeyFormat, tokenKey)

	tokenValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(constants.ErrorToGetTokenFromRedis, commonkeys.TokenKey, cacheKey, commonkeys.Error, err)
		return domain.Token{}, err
	}

	if tokenValue == "" {
		span.SetStatus(codes.Ok, sharederrors.ErrTokenNotFound)
		return domain.Token{}, nil
	}

	tokenDomain := domain.Token{
		Key:   tokenKey,
		Value: tokenValue,
	}

	span.SetStatus(codes.Ok, constants.TokenRetrievedSuccessfully)
	return tokenDomain, nil
}
