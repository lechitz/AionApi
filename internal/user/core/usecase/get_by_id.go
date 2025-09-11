package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetByID retrieves a user by the provided userID.
func (s *Service) GetByID(ctx context.Context, userID uint64) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetSelf)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetSelf),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToGetSelf))
		s.logger.ErrorwCtx(ctx, ErrorToGetSelf, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return domain.User{}, err
	}

	span.SetAttributes(attribute.String(commonkeys.Status, commonkeys.StatusSuccess))
	s.logger.InfowCtx(ctx, SuccessUserRetrieved, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, nil
}
