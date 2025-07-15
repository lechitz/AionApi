package user

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

// GetUserByUsername retrieves a user by their username from the database. Returns the user or an error if the operation fails.
//
//nolint:dupl // TODO: Refactor duplicated tracing logic
func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetUserByUsername)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetUserByUsername),
		attribute.String(commonkeys.Username, username),
	)

	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByUsername))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByUsername, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}
