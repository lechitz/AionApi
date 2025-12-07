// Package usecase contains the business logic for the Chat context.
// It orchestrates input commands, applies validation and domain rules,
// and delegates communication to the Aion-Chat service while handling
// observability and logging concerns.
package usecase

import (
	"github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/chat/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// ChatService provides operations for processing chat messages using the Aion-Chat AI service.
type ChatService struct {
	aionChatClient output.AionChatClient
	logger         logger.ContextLogger
}

// NewService creates and returns a new instance of ChatService with the given client and logger dependencies.
func NewService(client output.AionChatClient, log logger.ContextLogger) input.ChatService {
	return &ChatService{
		aionChatClient: client,
		logger:         log,
	}
}
