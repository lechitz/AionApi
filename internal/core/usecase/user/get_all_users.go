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

// GetAllUsers retrieves all users from the system. Returns a slice of User or an error if the operation fails.
func (s *Service) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetAllUsers)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetAllUsers),
	)

	users, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetAllUsers))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetAllUsers, commonkeys.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.Int(commonkeys.UsersCount, len(users)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUsersRetrieved, commonkeys.Users, strconv.Itoa(len(users)))
	return users, nil
}
