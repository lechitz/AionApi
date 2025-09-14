// Package usecase (user) contains use cases for managing users in the system.
package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
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
		return domain.User{}, err
	}

	hashed, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, StatusHashPasswordFailed))
		s.logger.ErrorwCtx(ctx, ErrorToHashPassword, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	user := domain.User{
		Name:     cmd.Name,
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: hashed,
		Roles:    UserRoles,
	}

	userDomain, err := s.userRepository.Create(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, StatusDBErrorCreateUser))
		s.logger.ErrorwCtx(ctx, ErrorToCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, err
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

// normalizeCreateCmd trims whitespace and normalizes a case for user creation input fields.
func normalizeCreateCmd(c *input.CreateUserCommand) {
	c.Name = strings.TrimSpace(c.Name)
	c.Username = strings.TrimSpace(c.Username)
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
}
