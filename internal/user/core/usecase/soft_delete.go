// Package usecase user contains use cases for managing users in the system.
package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// SoftDeleteUser performs a soft delete operation on a user identified by userID and deletes associated tokens. Returns an error if the operation fails.
func (s *Service) SoftDeleteUser(ctx context.Context, userID uint64) error {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanSoftDeleteUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanSoftDeleteUser),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	if err := s.authStore.Delete(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgDeleteToken))
		s.logger.ErrorwCtx(ctx, sharederrors.ErrMsgDeleteToken, commonkeys.Error, err.Error())
		return err
	}

	if err := s.userRepository.SoftDelete(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToSoftDeleteUser))
		s.logger.ErrorwCtx(ctx, ErrorToSoftDeleteUser, commonkeys.Error, err.Error())
		return err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, SuccessUserSoftDeleted, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return nil
}
