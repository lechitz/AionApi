// Package usecase (user) contains use cases for managing users in the system.
package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// ListAll retrieves all users.
func (s *Service) ListAll(ctx context.Context) ([]domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetAllUsers)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetAllUsers),
	)

	users, err := s.userRepository.ListAll(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToGetSelf))
		s.logger.ErrorwCtx(ctx, ErrorToGetSelf, commonkeys.Error, err.Error())
		return []domain.User{}, err
	}

	span.SetAttributes(attribute.String(commonkeys.Status, commonkeys.StatusSuccess))
	s.logger.InfowCtx(ctx, SuccessUserRetrieved)
	return users, nil
}
