package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Save persists the token for the given userID with a TTL.
func (s *Store) Save(ctx context.Context, token domain.Auth, expiration time.Duration) error {
	tr := otel.Tracer(SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, SpanNameTokenSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
		attribute.String(commonkeys.TokenKey, strconv.FormatUint(token.Key, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TokenUserKeyFormat, token.Key, token.Type)

	// if no expiration passed, fall back to default
	if expiration <= 0 {
		expiration = TokenExpirationDefault
	}

	if err := s.cache.Set(ctx, cacheKey, token.Token, expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveTokenToRedis, commonkeys.TokenKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TokenSavedSuccessfully)
	return nil
}
