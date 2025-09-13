package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create is the handler for creating a new category.
func (h *controller) Create(ctx context.Context, input model.CreateCategoryInput, userID uint64) (*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreate),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	category := toDomainCreate(input, userID)
	categoryDomain, err := h.CategoryService.Create(ctx, category)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, StatusCreated)
	return toModelOut(categoryDomain), nil
}
