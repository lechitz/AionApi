// Package usecase contains the business logic for the Chat context.
// It orchestrates input commands, applies validation and domain rules,
// and delegates communication to the Aion-Chat service while handling
// observability and logging concerns.
package usecase

import (
	auditinput "github.com/lechitz/AionApi/internal/audit/core/ports/input"
	"github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/chat/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// ChatService provides operations for processing chat messages using the Aion-Chat AI service.
type ChatService struct {
	aionChatClient   output.AionChatClient
	chatHistoryRepo  output.ChatHistoryRepository
	chatHistoryCache output.ChatHistoryCache // Redis cache for fast history access
	auditService     auditinput.Service
	logger           logger.ContextLogger
}

// NewService creates and returns a new instance of ChatService with the given client and logger dependencies.
func NewService(
	client output.AionChatClient,
	historyRepo output.ChatHistoryRepository,
	historyCache output.ChatHistoryCache,
	auditSvc auditinput.Service,
	log logger.ContextLogger,
) input.ChatService {
	return &ChatService{
		aionChatClient:   client,
		chatHistoryRepo:  historyRepo,
		chatHistoryCache: historyCache,
		auditService:     auditSvc,
		logger:           log,
	}
}
