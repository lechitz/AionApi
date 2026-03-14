// Package usecase_test contains tests for chat use cases.
package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetChatContext_Success(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(555)

	// Mock: HistoryRepo returns recent chats
	recentChats := []domain.ChatHistory{
		{
			ChatID:     201,
			UserID:     userID,
			Message:    TestHistoryQuestionWater,
			Response:   TestHistoryAnswerWater,
			TokensUsed: 20,
			CreatedAt:  time.Now().Add(-1 * time.Hour),
		},
		{
			ChatID:     202,
			UserID:     userID,
			Message:    TestHistoryQuestionReminder,
			Response:   TestHistoryAnswerReminder,
			TokensUsed: 15,
			CreatedAt:  time.Now().Add(-30 * time.Minute),
		},
	}
	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, 5).
		Return(recentChats, nil)

	// Execute
	context, err := suite.ChatService.GetChatContext(suite.Ctx, userID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, context)
	require.Len(t, context.RecentChats, 2)
	require.Equal(t, uint64(201), context.RecentChats[0].ChatID)
	require.Equal(t, TestHistoryQuestionWater, context.RecentChats[0].Message)
}

func TestGetChatContext_EmptyHistory(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(666)

	// Mock: HistoryRepo returns empty (new user)
	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, 5).
		Return([]domain.ChatHistory{}, nil)

	// Execute
	context, err := suite.ChatService.GetChatContext(suite.Ctx, userID)

	// Assert: Should succeed with empty context
	require.NoError(t, err)
	require.NotNil(t, context)
	require.Empty(t, context.RecentChats)
}

func TestGetChatContext_RepositoryError_NonCritical(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(777)

	// Mock: HistoryRepo fails (non-critical error)
	repoError := errors.New(TestErrorDatabaseTimeout)
	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, 5).
		Return(nil, repoError)

	// Execute
	context, err := suite.ChatService.GetChatContext(suite.Ctx, userID)

	// Assert: Should succeed with empty RecentChats (error is non-blocking)
	require.NoError(t, err)
	require.NotNil(t, context)
	require.Empty(t, context.RecentChats)
}
