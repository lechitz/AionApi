// Package setup provides test fixtures and mocks for chat use cases.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/chat/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// ChatServiceTestSuite holds all dependencies needed to test the chat service.
type ChatServiceTestSuite struct {
	Ctx            context.Context
	Ctrl           *gomock.Controller
	ChatService    input.ChatService
	AionChatClient *mocks.MockAionChatClient
	HistoryRepo    *mocks.MockChatHistoryRepository
	HistoryCache   *mocks.MockChatHistoryCache
	AuditService   *mocks.MockService
	Logger         *mocks.MockContextLogger
}

// ChatServiceTest initializes a test suite for the chat service.
func ChatServiceTest(t *testing.T) *ChatServiceTestSuite {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	aionChatClient := mocks.NewMockAionChatClient(ctrl)
	historyRepo := mocks.NewMockChatHistoryRepository(ctrl)
	historyCache := mocks.NewMockChatHistoryCache(ctrl)
	auditService := mocks.NewMockService(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)

	ExpectLoggerDefaultBehavior(logger)

	service := usecase.NewService(
		aionChatClient,
		historyRepo,
		historyCache,
		auditService,
		logger,
	)

	return &ChatServiceTestSuite{
		Ctx:            ctx,
		Ctrl:           ctrl,
		ChatService:    service,
		AionChatClient: aionChatClient,
		HistoryRepo:    historyRepo,
		HistoryCache:   historyCache,
		AuditService:   auditService,
		Logger:         logger,
	}
}
