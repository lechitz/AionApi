// Package usecase user contains use cases for managing users in the system.
package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		span.RecordError(ErrNoFieldsToUpdate)
		span.SetStatus(codes.Error, ErrorNoFieldsToUpdate)
		s.logger.ErrorwCtx(ctx, ErrorNoFieldsToUpdate, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return domain.User{}, ErrNoFieldsToUpdate
	}

	fields := buildUpdateFields(cmd)

	updatedUser, err := s.userRepository.Update(ctx, userID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToUpdateUser))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToUpdateUser, commonkeys.Error, err.Error())
		return domain.User{}, fmt.Errorf("%w: %w", ErrUpdateUser, err)
	}

	span.AddEvent(SpanEventInvalidateCache)
	if err := s.userCache.DeleteUser(ctx, updatedUser.ID, updatedUser.Username, updatedUser.Email); err != nil {
		s.logger.WarnwCtx(ctx, WarnFailedToInvalidateUserCache,
			commonkeys.UserID, updatedUser.ID,
			commonkeys.Username, updatedUser.Username,
			commonkeys.Email, updatedUser.Email,
			commonkeys.Error, err,
		)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessUserUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10), commonkeys.UserUpdatedFields, updatedUser)

	return updatedUser, nil
}

// buildUpdateFields constructs a map of fields to update from the command.
// Only non-nil fields from the command are included in the update.
// Username and email are normalized to lowercase to maintain consistency.
func buildUpdateFields(cmd input.UpdateUserCommand) map[string]interface{} {
	fields := map[string]interface{}{
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	if cmd.Name != nil {
		fields[commonkeys.Name] = strings.TrimSpace(*cmd.Name)
	}
	if cmd.Username != nil {
		fields[commonkeys.Username] = strings.ToLower(strings.TrimSpace(*cmd.Username))
	}
	if cmd.Email != nil {
		fields[commonkeys.Email] = strings.ToLower(strings.TrimSpace(*cmd.Email))
	}
	if cmd.Locale != nil {
		fields[commonkeys.Locale] = *cmd.Locale
	}
	if cmd.Timezone != nil {
		fields[commonkeys.Timezone] = *cmd.Timezone
	}
	if cmd.Location != nil {
		fields[commonkeys.Location] = *cmd.Location
	}
	if cmd.Bio != nil {
		fields[commonkeys.Bio] = *cmd.Bio
	}
	if cmd.AvatarURL != nil {
		fields[commonkeys.AvatarURL] = *cmd.AvatarURL
	}

	return fields
}
