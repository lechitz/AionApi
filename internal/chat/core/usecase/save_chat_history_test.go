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

func TestSaveChatHistory_Success(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	message := "What is AI?"
	response := "AI stands for Artificial Intelligence..."
	tokensUsed := 100
	functionCalls := map[string]string{
		"search_knowledge": `{"query": "artificial intelligence"}`,
	}

	savedHistory := domain.ChatHistory{
		ChatID:        789,
		UserID:        userID,
		Message:       message,
		Response:      response,
		TokensUsed:    tokensUsed,
		FunctionCalls: functionCalls,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Mock repository save
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ interface{}, history domain.ChatHistory) (domain.ChatHistory, error) {
			require.Equal(t, userID, history.UserID)
			require.Equal(t, message, history.Message)
			require.Equal(t, response, history.Response)
			require.Equal(t, tokensUsed, history.TokensUsed)
			require.Equal(t, functionCalls, history.FunctionCalls)
			return savedHistory, nil
		})

	// Mock cache add
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, savedHistory).
		Return(nil)

	err := suite.ChatService.SaveChatHistory(suite.Ctx, userID, message, response, tokensUsed, functionCalls)

	require.NoError(t, err)
}

func TestSaveChatHistory_RepositoryError(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	message := "Test"
	response := "Response"
	tokensUsed := 10
	functionCalls := map[string]string{}

	// Mock repository returning error
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, errors.New("database connection failed"))

	err := suite.ChatService.SaveChatHistory(suite.Ctx, userID, message, response, tokensUsed, functionCalls)

	require.Error(t, err)
	require.ErrorContains(t, err, "database connection failed")
}

func TestSaveChatHistory_CacheError_StillSucceeds(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	message := "Test"
	response := "Response"
	tokensUsed := 10
	functionCalls := map[string]string{}

	savedHistory := domain.ChatHistory{
		ChatID:        123,
		UserID:        userID,
		Message:       message,
		Response:      response,
		TokensUsed:    tokensUsed,
		FunctionCalls: functionCalls,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(savedHistory, nil)

	// Mock cache returning error (should not fail the operation)
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, savedHistory).
		Return(errors.New("redis unavailable"))

	err := suite.ChatService.SaveChatHistory(suite.Ctx, userID, message, response, tokensUsed, functionCalls)

	// Should succeed despite cache error
	require.NoError(t, err)
}

func TestSaveChatHistory_WithEmptyFunctionCalls(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	message := "Hello"
	response := "Hi!"
	tokensUsed := 5
	functionCalls := map[string]string{} // Empty map

	savedHistory := domain.ChatHistory{
		ChatID:        111,
		UserID:        userID,
		Message:       message,
		Response:      response,
		TokensUsed:    tokensUsed,
		FunctionCalls: functionCalls,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(savedHistory, nil)

	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, savedHistory).
		Return(nil)

	err := suite.ChatService.SaveChatHistory(suite.Ctx, userID, message, response, tokensUsed, functionCalls)

	require.NoError(t, err)
}

func TestSaveChatHistory_WithLongMessage(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	// Simulate a very long message (e.g., 1000 chars)
	message := string(make([]byte, 1000))
	for i := range message {
		message = message[:i] + "a" + message[i+1:]
	}
	response := "Processed long message"
	tokensUsed := 250
	functionCalls := map[string]string{}

	savedHistory := domain.ChatHistory{
		ChatID:        222,
		UserID:        userID,
		Message:       message,
		Response:      response,
		TokensUsed:    tokensUsed,
		FunctionCalls: functionCalls,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ interface{}, history domain.ChatHistory) (domain.ChatHistory, error) {
			require.Len(t, history.Message, 1000)
			return savedHistory, nil
		})

	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, savedHistory).
		Return(nil)

	err := suite.ChatService.SaveChatHistory(suite.Ctx, userID, message, response, tokensUsed, functionCalls)

	require.NoError(t, err)
}
