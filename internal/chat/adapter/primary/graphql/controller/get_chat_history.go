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

// GetChatHistory retrieves chat history for a user with pagination support.
func (c *controller) GetChatHistory(ctx context.Context, userID uint64, limit, offset int) ([]*model.ChatMessage, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetChatHistory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrLimit, limit),
		attribute.Int(AttrOffset, offset),
	)

	// Delegate to the service (use case)
	histories, err := c.ChatService.GetChatHistory(ctx, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgFetchError)
		c.Logger.ErrorwCtx(
			ctx,
			MsgFetchError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			AttrLimit, limit,
			AttrOffset, offset,
		)
		return nil, err
	}

	// Convert to GraphQL model
	result := toModelOutSlice(histories)

	span.SetAttributes(attribute.Int(AttrCount, len(result)))
	span.SetStatus(codes.Ok, StatusFetched)

	c.Logger.InfowCtx(
		ctx,
		MsgFetched,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		AttrCount, len(result),
	)

	return result, nil
}
