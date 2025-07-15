// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CreateUser creates a new user with the given data and password, ensuring validations and unique constraints are met. Returns the created user or an error.
func (s *Service) CreateUser(ctx context.Context, newUser domain.User) (domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanCreateUser)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanCreateUser),
		attribute.String(commonkeys.Username, newUser.Username),
		attribute.String(commonkeys.Email, newUser.Email),
	)

	s.normalizeUserData(&newUser)

	//TODO: talvez esse método abaixo possa morrer. Parar pra ver depois.

	if err := s.validateCreateUserRequired(newUser); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusValidationFailed))
		s.logger.ErrorwCtx(ctx, constants.ErrorToValidateCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, sharederrors.NewValidationError(commonkeys.User, err.Error())
	}

	//TODO: criar um novo método que valida Username e Email de uma só vez: ExistsByUsernameOrEmail

	existingByUsername, err := s.userRepository.GetUserByUsername(ctx, newUser.Username)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusDBErrorUsername))
		s.logger.ErrorwCtx(ctx, constants.DBErrorCheckingUsername, commonkeys.Error, err.Error())
		return domain.User{}, err
	}
	if existingByUsername.ID != 0 {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusUsernameExists))
		return domain.User{}, sharederrors.NewValidationError(commonkeys.Username, sharederrors.ErrUsernameInUse)
	}

	existingByEmail, err := s.userRepository.GetUserByEmail(ctx, newUser.Email)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusDBErrorEmail))
		s.logger.ErrorwCtx(ctx, constants.DBErrorCheckingEmail, commonkeys.Error, err.Error())
		return domain.User{}, err
	}
	if existingByEmail.ID != 0 {
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusEmailExists))
		return domain.User{}, sharederrors.NewValidationError(commonkeys.Email, sharederrors.ErrEmailInUse)
	}

	hashedPassword, err := s.hashStore.Hash(newUser.Password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusHashPasswordFailed))
		s.logger.ErrorwCtx(ctx, constants.ErrorToHashPassword, commonkeys.Error, err.Error())
		return domain.User{}, err
	}
	newUser.Password = hashedPassword

	userCreated, err := s.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.StatusDBErrorCreateUser))
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userCreated.ID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserCreated, commonkeys.UserID, strconv.FormatUint(userCreated.ID, 10))

	return userCreated, nil
}
