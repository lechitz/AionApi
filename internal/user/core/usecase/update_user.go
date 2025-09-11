// Package user contains use cases for managing users in the system.
package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, userID uint64, cmd input.UpdateUserCommand) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUpdateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateUser),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	if !cmd.HasUpdates() {
		span.RecordError(sharederrors.ErrNoFieldsToUpdate)
		span.SetStatus(codes.Error, ErrorNoFieldsToUpdate)
		s.logger.ErrorwCtx(ctx, ErrorNoFieldsToUpdate, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return domain.User{}, errors.New(ErrorNoFieldsToUpdate)
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
	fields[UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userRepository.Update(ctx, userID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToUpdateUser))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToUpdateUser, commonkeys.Error, err.Error())
		return domain.User{}, fmt.Errorf("%s: %w", ErrorToUpdateUser, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessUserUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10), commonkeys.UserUpdatedFields, updatedUser)

	return updatedUser, nil
}
