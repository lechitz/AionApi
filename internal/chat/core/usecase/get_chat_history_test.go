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

func TestGetChatHistory_Success(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	limit := 10
	offset := 0

	expectedHistories := []domain.ChatHistory{
		{
			ChatID:     1,
			UserID:     userID,
			Message:    "Hello",
			Response:   "Hi there!",
			TokensUsed: 5,
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			ChatID:     2,
			UserID:     userID,
			Message:    "What is AI?",
			Response:   "AI stands for...",
			TokensUsed: 50,
			CreatedAt:  time.Now().Add(-1 * time.Hour),
		},
	}

	suite.HistoryRepo.EXPECT().
		GetByUserID(gomock.Any(), userID, limit, offset).
		Return(expectedHistories, nil)

	histories, err := suite.ChatService.GetChatHistory(suite.Ctx, userID, limit, offset)

	require.NoError(t, err)
	require.Len(t, histories, 2)
	require.Equal(t, expectedHistories[0].ChatID, histories[0].ChatID)
	require.Equal(t, expectedHistories[1].ChatID, histories[1].ChatID)
}

func TestGetChatHistory_WithPagination(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	limit := 5
	offset := 10

	expectedHistories := []domain.ChatHistory{
		{ChatID: 11, UserID: userID, Message: "Page 2 message 1", Response: "Response 1"},
		{ChatID: 12, UserID: userID, Message: "Page 2 message 2", Response: "Response 2"},
	}

	suite.HistoryRepo.EXPECT().
		GetByUserID(gomock.Any(), userID, limit, offset).
		Return(expectedHistories, nil)

	histories, err := suite.ChatService.GetChatHistory(suite.Ctx, userID, limit, offset)

	require.NoError(t, err)
	require.Len(t, histories, 2)
}

func TestGetChatHistory_EmptyResult(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	limit := 10
	offset := 0

	suite.HistoryRepo.EXPECT().
		GetByUserID(gomock.Any(), userID, limit, offset).
		Return([]domain.ChatHistory{}, nil)

	histories, err := suite.ChatService.GetChatHistory(suite.Ctx, userID, limit, offset)

	require.NoError(t, err)
	require.Empty(t, histories)
}

func TestGetChatHistory_RepositoryError(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	limit := 10
	offset := 0

	suite.HistoryRepo.EXPECT().
		GetByUserID(gomock.Any(), userID, limit, offset).
		Return(nil, errors.New("database error"))

	histories, err := suite.ChatService.GetChatHistory(suite.Ctx, userID, limit, offset)

	require.Error(t, err)
	require.Nil(t, histories)
	require.ErrorContains(t, err, "database error")
}

func TestGetLatestChatHistory_Success(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	limit := 5

	expectedHistories := []domain.ChatHistory{
		{
			ChatID:     100,
			UserID:     userID,
			Message:    "Recent message",
			Response:   "Recent response",
			TokensUsed: 20,
			CreatedAt:  time.Now().Add(-5 * time.Minute),
		},
	}

	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, limit).
		Return(expectedHistories, nil)

	histories, err := suite.ChatService.GetLatestChatHistory(suite.Ctx, userID, limit)

	require.NoError(t, err)
	require.Len(t, histories, 1)
	require.Equal(t, uint64(100), histories[0].ChatID)
	require.Equal(t, "Recent message", histories[0].Message)
}

func TestGetLatestChatHistory_Empty(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(789)
	limit := 5

	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, limit).
		Return([]domain.ChatHistory{}, nil)

	histories, err := suite.ChatService.GetLatestChatHistory(suite.Ctx, userID, limit)

	require.NoError(t, err)
	require.Empty(t, histories)
}

func TestGetLatestChatHistory_RepositoryError(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	limit := 5

	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, limit).
		Return(nil, errors.New("connection timeout"))

	histories, err := suite.ChatService.GetLatestChatHistory(suite.Ctx, userID, limit)

	require.Error(t, err)
	require.Nil(t, histories)
	require.ErrorContains(t, err, "connection timeout")
}

func TestGetLatestChatHistory_LargeLimit(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	limit := 100 // Large limit

	// Create 100 history items
	expectedHistories := make([]domain.ChatHistory, 100)
	for i, chatID := 0, uint64(1); i < len(expectedHistories); i, chatID = i+1, chatID+1 {
		expectedHistories[i] = domain.ChatHistory{
			ChatID:     chatID,
			UserID:     userID,
			Message:    "Message " + string(rune(i)),
			Response:   "Response " + string(rune(i)),
			TokensUsed: 10,
		}
	}

	suite.HistoryRepo.EXPECT().
		GetLatest(gomock.Any(), userID, limit).
		Return(expectedHistories, nil)

	histories, err := suite.ChatService.GetLatestChatHistory(suite.Ctx, userID, limit)

	require.NoError(t, err)
	require.Len(t, histories, 100)
}
