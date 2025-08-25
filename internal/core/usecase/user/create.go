// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Create handles user creation by normalizing input data, checking for uniqueness conflicts,
// hashing the password, and persisting the user in the repository.
// It returns the created user or an error if the operation fails.
func (s *Service) Create(ctx context.Context, cmd input.CreateUserCommand) (domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanCreateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanCreateUser),
		attribute.String(commonkeys.Username, cmd.Username),
		attribute.String(commonkeys.Email, cmd.Email),
	)

	normalizeCreateCmd(&cmd)

	conflict, err := s.userRepository.CheckUniqueness(ctx, cmd.Username, cmd.Email)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusDBErrorUsernameOrEmail))
		return domain.User{}, err
	}
	if conflict.UsernameTaken {
		valErr := sharederrors.NewValidationError(commonkeys.Username, sharederrors.ErrUsernameInUse)
		span.RecordError(valErr)
		return domain.User{}, valErr
	}
	if conflict.EmailTaken {
		valErr := sharederrors.NewValidationError(commonkeys.Email, sharederrors.ErrEmailInUse)
		span.RecordError(valErr)
		return domain.User{}, valErr
	}

	hashed, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusHashPasswordFailed))
		s.logger.ErrorwCtx(ctx, constants.ErrorToHashPassword, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	toSave := domain.User{
		Name:     cmd.Name,
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: hashed,
	}

	userDomain, err := s.userRepository.Create(ctx, toSave)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusDBErrorCreateUser))
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userDomain.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserCreated, commonkeys.UserID, strconv.FormatUint(userDomain.ID, 10))

	return userDomain, nil
}

// normalizeCreateCmd trims whitespace and normalizes a case for user creation input fields.
func normalizeCreateCmd(c *input.CreateUserCommand) {
	c.Name = strings.TrimSpace(c.Name)
	c.Username = strings.TrimSpace(c.Username)
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
}
