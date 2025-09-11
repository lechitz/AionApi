// Package usecase (auth) contains use cases for authenticating users and generating tokens.
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
	ctx, span := tracer.Start(ctx, SpanRevokeToken,
		trace.WithAttributes(
			attribute.String(commonkeys.Operation, SpanRevokeToken),
			attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		),
	)
	defer span.End()

	span.AddEvent(EventRevokeToken)
	if err := s.authStore.Delete(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToDeleteToken)
		s.logger.ErrorwCtx(ctx, ErrorToDeleteToken, commonkeys.Error, err.Error())
		return fmt.Errorf("%s: %w", ErrorToDeleteToken, err)
	}

	span.AddEvent(EventTokenRevoked)
	span.SetStatus(codes.Ok, SuccessUserLoggedOut)

	s.logger.InfowCtx(ctx, SuccessUserLoggedOut,
		commonkeys.TokenKey, strconv.FormatUint(userID, 10),
	)
	return nil
}
