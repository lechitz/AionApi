// Package handler provides HTTP handlers for the chat module.
package handler

import (
	"github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// @tag.name Chat
// @tag.description AI Chat endpoints for conversational queries

// Handler wires chat use cases to HTTP handlers.
type Handler struct {
	Service input.ChatService
	Logger  logger.ContextLogger
	Config  *config.Config
}

// New creates a new Handler instance with its dependencies wired.
func New(chatService input.ChatService, cfg *config.Config, log logger.ContextLogger) *Handler {
	return &Handler{
		Service: chatService,
		Config:  cfg,
		Logger:  log,
	}
}
