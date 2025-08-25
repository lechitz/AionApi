// internal/core/usecase/auth/logout.go
package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Logout revokes a user's authentication token.
func (s *Service) Logout(ctx context.Context, userID uint64) error {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanLogout)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanLogout),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	span.AddEvent(constants.EventRevokeToken)
	if err := s.tokenStore.Delete(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToRevokeToken)
		s.logger.ErrorwCtx(ctx, constants.ErrorToRevokeToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return fmt.Errorf("%s: %w", constants.ErrorToRevokeToken, err)
	}

	span.AddEvent(constants.EventLogoutSuccess)
	span.SetStatus(codes.Ok, constants.SuccessUserLoggedOut)

	s.logger.InfowCtx(ctx, constants.SuccessUserLoggedOut, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return nil
}
