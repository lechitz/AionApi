// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateUser),
		attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)),
	)

	updateFields := buildUserUpdateFields(user)

	if len(updateFields) == 0 {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.NoFieldsToUpdate))
		err := errors.New(constants.ErrorNoFieldsToUpdate)
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorNoFieldsToUpdate, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	updateFields[constants.UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, updateFields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToUpdateUser))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, constants.ErrorToUpdateUser, commonkeys.Error, err.Error())
		return domain.User{}, fmt.Errorf("%s: %w", constants.ErrorToUpdateUser, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10), commonkeys.UserUpdatedFields, updateFields)

	return updatedUser, nil
}
