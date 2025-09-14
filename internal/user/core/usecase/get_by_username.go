package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetUserByUsername retrieves a user by their username from the database. Returns the user or an error if the operation fails.
func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetUserByUsername)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetUserByUsername),
		attribute.String(commonkeys.Username, username),
	)

	user, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToGetUserByUsername))
		s.logger.ErrorwCtx(ctx, ErrorToGetUserByUsername, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
	)
	s.logger.InfowCtx(ctx, SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}
