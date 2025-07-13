// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetUserByID retrieves a user by their unique ID from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByID(ctx context.Context, userID uint64) (domain.UserDomain, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetUserByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetUserByID),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	user, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByID))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByID, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}

// GetUserByEmail retrieves a user by their email address from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetUserByEmail)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetUserByEmail),
		attribute.String(commonkeys.Email, email),
	)

	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByEmail))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByEmail, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}

// GetUserByUsername retrieves a user by their username from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetUserByUsername)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetUserByUsername),
		attribute.String(commonkeys.Username, username),
	)

	user, err := s.userStore.GetUserByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByUsername))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByUsername, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}

// GetAllUsers retrieves all users from the system. Returns a slice of UserDomain or an error if the operation fails.
func (s *Service) GetAllUsers(ctx context.Context) ([]domain.UserDomain, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetAllUsers)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetAllUsers),
	)

	users, err := s.userStore.GetAllUsers(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetAllUsers))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetAllUsers, commonkeys.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.StatusSuccess),
		attribute.Int(commonkeys.UsersCount, len(users)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUsersRetrieved, commonkeys.Users, strconv.Itoa(len(users)))
	return users, nil
}
