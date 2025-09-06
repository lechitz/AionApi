package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Delete removes the token for the given user.
func (s *Store) Delete(ctx context.Context, tokenKey uint64) error {
	tr := otel.Tracer(SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, SpanNameTokenDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(tokenKey, 10)),
		attribute.String(commonkeys.Operation, commonkeys.OperationDelete),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TokenUserKeyFormat, tokenKey)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteTokenFromRedis, commonkeys.TokenKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TokenDeletedSuccessfully)
	return nil
}
