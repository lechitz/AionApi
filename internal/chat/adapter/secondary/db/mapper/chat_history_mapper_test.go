package mapper_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lechitz/aion-api/internal/chat/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/chat/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestChatHistoryToDBAndFromDB(t *testing.T) {
	now := time.Now().UTC()
	sid := uuid.New()
	execTime := 45
	deleted := now.Add(-time.Hour)

	src := domain.ChatHistory{
		ChatID:          1,
		UserID:          2,
		Message:         "hello",
		Response:        "world",
		TokensUsed:      10,
		FunctionCalls:   map[string]string{"a": "b"},
		CreatedAt:       now,
		UpdatedAt:       now,
		DeletedAt:       &deleted,
		SessionID:       &sid,
		ExecutionTimeMs: &execTime,
		ToolCount:       2,
		ErrorCount:      1,
	}

	db := mapper.ChatHistoryToDB(src)
	require.NotEmpty(t, db.FunctionCalls)
	require.True(t, db.DeletedAt.Valid)

	back := mapper.ChatHistoryFromDB(db)
	require.Equal(t, src.ChatID, back.ChatID)
	require.Equal(t, src.UserID, back.UserID)
	require.Equal(t, src.FunctionCalls, back.FunctionCalls)
	require.NotNil(t, back.DeletedAt)
	require.Equal(t, deleted, *back.DeletedAt)
}

func TestChatHistoryFromDB_InvalidFunctionCallsAndSliceHelper(t *testing.T) {
	now := time.Now().UTC()
	db := model.ChatHistoryDB{
		ChatID:        3,
		UserID:        9,
		Message:       "m",
		Response:      "r",
		FunctionCalls: []byte("{invalid"),
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     gorm.DeletedAt{},
	}

	out := mapper.ChatHistoryFromDB(db)
	require.Nil(t, out.FunctionCalls)
	require.Nil(t, out.DeletedAt)

	slice := mapper.ChatHistoriesFromDB([]model.ChatHistoryDB{db})
	require.Len(t, slice, 1)
	require.Equal(t, uint64(3), slice[0].ChatID)
}
