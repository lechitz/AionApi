// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// SoftDeleteUser performs a soft delete operation on a user identified by userID and deletes associated tokens. Returns an error if the operation fails.
func (s *Service) SoftDeleteUser(ctx context.Context, userID uint64) error {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanSoftDeleteUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanSoftDeleteUser),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	if err := s.userStore.SoftDeleteUser(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToSoftDeleteUser))
		s.logger.ErrorwCtx(ctx, constants.ErrorToSoftDeleteUser, commonkeys.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.tokenService.Delete(ctx, tokenDomain); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrorToDeleteToken))
		s.logger.ErrorwCtx(ctx, sharederrors.ErrorToDeleteToken, commonkeys.Error, err.Error())
		return err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserSoftDeleted, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return nil
}
