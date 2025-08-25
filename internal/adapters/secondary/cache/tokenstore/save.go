package tokenstore

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/tokenstore/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Save persists the token for the given userID.
func (s *Store) Save(ctx context.Context, token domain.Token) error {
	tr := otel.Tracer(constants.SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, constants.SpanNameTokenSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, commonkeys.OperationSave),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
		attribute.String(commonkeys.TokenKey, strconv.FormatUint(token.Key, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(constants.TokenUserKeyFormat, token.Key)

	if err := s.cache.Set(ctx, cacheKey, token.Value, constants.TokenExpirationDefault); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(constants.ErrorToSaveTokenToRedis, commonkeys.TokenKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, constants.TokenSavedSuccessfully)
	return nil
}
