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

// GetUserByID retrieves a user by their unique ID from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByID(ctx context.Context, userID uint64) (domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetUserByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetUserByID),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByID))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByID, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}
