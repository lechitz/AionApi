package category

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// ListAll retrieves all categories associated with a specific user ID using the repository.
// Returns a list of categories or an error in case of failure.
func (s *Service) ListAll(ctx context.Context, userID uint64) ([]domain.Category, error) {
	tr := otel.Tracer(constants.TracerName)
	ctx, span := tr.Start(ctx, constants.SpanListAllCategories)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanListAllCategories),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	span.AddEvent(constants.EventRepositoryListAll)
	categories, err := s.CategoryRepository.ListAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.FailedToGetAllCategories)
		s.Logger.ErrorwCtx(ctx, constants.FailedToGetAllCategories, commonkeys.Error, err)
		return nil, err
	}

	span.AddEvent(constants.EventSuccess)
	span.SetStatus(codes.Ok, "retrieved")
	return categories, nil
}
