// Package usecase (user) contains use cases for managing users in the system.
package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Create handles user creation by normalizing input data, checking for uniqueness conflicts,
// hashing the password, and persisting the user in the repository.
// It returns the created user or an error if the operation fails.
func (s *Service) Create(ctx context.Context, cmd input.CreateUserCommand) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanCreateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreateUser),
		attribute.String(commonkeys.Username, cmd.Username),
		attribute.String(commonkeys.Email, cmd.Email),
	)

	normalizeCreateCmd(&cmd)

	conflict, err := s.userRepository.CheckUniqueness(ctx, cmd.Username, cmd.Email)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, StatusDBErrorUsernameOrEmail))
		return domain.User{}, err
	}

	if err := validateUniquenessConflicts(ctx, span, conflict.UsernameTaken, conflict.EmailTaken); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, StatusUsernameOrEmailInUse))
		s.logger.WarnwCtx(ctx, WarnUsernameOrEmailInUse,
			commonkeys.Username, cmd.Username,
			commonkeys.Email, cmd.Email,
		)
		return domain.User{}, err
	}

	hashed, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, StatusHashPasswordFailed))
		s.logger.ErrorwCtx(ctx, ErrorToHashPassword, commonkeys.Error, err.Error())
		return domain.User{}, ErrHashPassword
	}

	user := domain.User{
		Name:      cmd.Name,
		Username:  cmd.Username,
		Email:     cmd.Email,
		Password:  hashed,
		Locale:    cmd.Locale,
		Timezone:  cmd.Timezone,
		Location:  cmd.Location,
		Bio:       cmd.Bio,
		AvatarURL: cmd.AvatarURL,
	}

	userDomain, err := s.userRepository.Create(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, StatusDBErrorCreateUser))
		s.logger.ErrorwCtx(ctx, ErrorToCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, ErrCreateUser
	}

	span.AddEvent(SpanEventSaveToCache)
	if err := s.userCache.SaveUser(ctx, userDomain, 0); err != nil {
		s.logger.WarnwCtx(ctx, WarnFailedToSaveUserToCache,
			commonkeys.UserID, userDomain.ID,
			commonkeys.Error, err,
		)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userDomain.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessUserCreated, commonkeys.UserID, strconv.FormatUint(userDomain.ID, 10))

	return userDomain, nil
}

// validateUniquenessConflicts checks if the username or email are already taken and returns an error if so.
func validateUniquenessConflicts(_ context.Context, span trace.Span, usernameTaken, emailTaken bool) error {
	if usernameTaken {
		valErr := sharederrors.NewValidationError(commonkeys.Username, sharederrors.ErrUsernameInUse)
		span.RecordError(valErr)
		return valErr
	}
	if emailTaken {
		valErr := sharederrors.NewValidationError(commonkeys.Email, sharederrors.ErrEmailInUse)
		span.RecordError(valErr)
		return valErr
	}
	return nil
}

// normalizeCreateCmd trims whitespace and normalizes case for user creation input fields.
// Username and email are converted to lowercase to prevent case-sensitive duplicates.
func normalizeCreateCmd(c *input.CreateUserCommand) {
	c.Name = strings.TrimSpace(c.Name)
	c.Username = strings.ToLower(strings.TrimSpace(c.Username))
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
}
