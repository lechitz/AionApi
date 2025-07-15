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

// GetUserByEmail retrieves a user by their email address from the database. Returns the user or an error if the operation fails.
//
//nolint:dupl // TODO: Refactor duplicated tracing logic
func (s *Service) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetUserByEmail)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetUserByEmail),
		attribute.String(commonkeys.Email, email),
	)

	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetUserByEmail))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByEmail, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}
