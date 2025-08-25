package token

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// Revoke deletes the user's token from the store. It is idempotent from the use case perspective.
func (s *Service) Revoke(ctx context.Context, tokenKey uint64) error {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(
		ctx,
		constants.SpanRevokeToken,
		trace.WithAttributes(
			attribute.String(commonkeys.Operation, constants.SpanRevokeToken),
			attribute.String(commonkeys.TokenKey, strconv.FormatUint(tokenKey, 10)),
		),
	)
	defer span.End()

	span.AddEvent(constants.EventRevokeToken)
	if err := s.tokenStore.Delete(ctx, tokenKey); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToDeleteToken)
		s.logger.ErrorwCtx(ctx, constants.ErrorToDeleteToken, commonkeys.Error, err.Error())
		return fmt.Errorf("%s: %w", constants.ErrorToDeleteToken, err)
	}

	span.AddEvent(constants.EventTokenRevoked)
	span.SetStatus(codes.Ok, constants.SuccessTokenDeleted)

	s.logger.InfowCtx(ctx, constants.SuccessTokenDeleted,
		commonkeys.TokenKey, strconv.FormatUint(tokenKey, 10),
	)
	return nil
}
