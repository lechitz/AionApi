// Package user contains use cases for managing users in the system.
package user

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// ListAll retrieves all users.
func (s *Service) ListAll(ctx context.Context) ([]domain.User, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanGetAllUsers)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetAllUsers),
	)

	users, err := s.userRepository.ListAll(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetSelf))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetSelf, commonkeys.Error, err.Error())
		return []domain.User{}, err
	}

	span.SetAttributes(attribute.String(commonkeys.Status, commonkeys.StatusSuccess))
	s.logger.InfowCtx(ctx, constants.SuccessUserRetrieved)
	return users, nil
}
