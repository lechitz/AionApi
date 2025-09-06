// internal/core/usecase/auth/logout.go
package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Logout revokes a user's authentication token.
func (s *Service) Logout(ctx context.Context, userID uint64) error {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(
		ctx,
		SpanRevokeToken,
		trace.WithAttributes(
			attribute.String(commonkeys.Operation, SpanRevokeToken),
			attribute.String(commonkeys.TokenKey, strconv.FormatUint(tokenKey, 10)),
		),
	)
	defer span.End()

	span.AddEvent(EventRevokeToken)
	if err := s.tokenStore.Delete(ctx, tokenKey); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToDeleteToken)
		s.logger.ErrorwCtx(ctx, ErrorToDeleteToken, commonkeys.Error, err.Error())
		return fmt.Errorf("%s: %w", ErrorToDeleteToken, err)
	}

	span.AddEvent(EventTokenRevoked)
	span.SetStatus(codes.Ok, SuccessTokenDeleted)

	s.logger.InfowCtx(ctx, SuccessTokenDeleted,
		commonkeys.TokenKey, strconv.FormatUint(tokenKey, 10),
	)
	return nil
}
