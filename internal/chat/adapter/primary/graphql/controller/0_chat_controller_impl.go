// Package controller contains GraphQL-facing controllers for the Chat context.
package controller

import (
	"context"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/chat/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// ChatController is the contract used by GraphQL resolvers.
type ChatController interface {
	GetChatHistory(ctx context.Context, userID uint64, limit, offset int) ([]*model.ChatMessage, error)
	GetChatContext(ctx context.Context, userID uint64) (*model.ChatContext, error)
}

// controller is the controller for the chat service.
type controller struct {
	ChatService input.ChatService
	Logger      logger.ContextLogger
}

// NewController wires dependencies and returns a ChatController.
func NewController(svc input.ChatService, logger logger.ContextLogger) ChatController {
	return &controller{
		ChatService: svc,
		Logger:      logger,
	}
}
