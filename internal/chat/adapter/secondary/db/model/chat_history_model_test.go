package model_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/chat/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestChatHistoryModelTableName(t *testing.T) {
	require.Equal(t, model.ChatHistoryTable, model.ChatHistoryDB{}.TableName())
}
