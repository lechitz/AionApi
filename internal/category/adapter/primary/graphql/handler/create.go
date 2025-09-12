package handler

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create mapeia o input GraphQL → domínio, delega ao service e retorna GraphQL model.
func (h *Handler) Create(ctx context.Context, input model.CreateCategoryInput, userID uint64) (*model.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreate),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	d := toDomainCreate(input, userID)
	created, err := h.CategoryService.Create(ctx, d)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, StatusCreated)
	return toModelOut(created), nil
}
