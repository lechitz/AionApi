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

// GetChatContext retrieves aggregated context data for AI assistance.
func (c *controller) GetChatContext(ctx context.Context, userID uint64) (*model.ChatContext, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "chat.controller.get_context")
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	// Delegate to the service (use case)
	contextData, err := c.ChatService.GetChatContext(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error fetching chat context")
		c.Logger.ErrorwCtx(
			ctx,
			"error fetching chat context",
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return nil, err
	}

	// Convert to GraphQL model
	result := &model.ChatContext{
		RecentChats:     toModelOutSlice(contextData.RecentChats),
		TotalRecords:    safeInt32(contextData.TotalRecords),
		TotalCategories: safeInt32(contextData.TotalCategories),
		TotalTags:       safeInt32(contextData.TotalTags),
	}

	span.SetAttributes(
		attribute.Int("recent_chats_count", len(result.RecentChats)),
	)
	span.SetStatus(codes.Ok, "chat context fetched")

	c.Logger.InfowCtx(
		ctx,
		"chat context fetched",
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"recent_chats", len(result.RecentChats),
	)

	return result, nil
}
