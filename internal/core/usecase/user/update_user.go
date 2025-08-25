// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, userID uint64, cmd input.UpdateUserCommand) (domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanUpdateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanUpdateUser),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	if !cmd.HasUpdates() {
		span.RecordError(sharederrors.ErrNoFieldsToUpdate)
		span.SetStatus(codes.Error, constants.ErrorNoFieldsToUpdate)
		s.logger.ErrorwCtx(ctx, constants.ErrorNoFieldsToUpdate, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return domain.User{}, errors.New(constants.ErrorNoFieldsToUpdate)
	}

	fields := map[string]interface{}{}
	if cmd.Name != nil {
		fields[commonkeys.Name] = *cmd.Name
	}
	if cmd.Username != nil {
		fields[commonkeys.Username] = *cmd.Username
	}
	if cmd.Email != nil {
		fields[commonkeys.Email] = *cmd.Email
	}
	fields[constants.UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userRepository.Update(ctx, userID, fields)
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
	s.logger.InfowCtx(ctx, constants.SuccessUserUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10), commonkeys.UserUpdatedFields, updatedUser)

	return updatedUser, nil
}
