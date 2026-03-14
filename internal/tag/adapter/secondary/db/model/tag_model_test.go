package model_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestTagModelTableName(t *testing.T) {
	require.Equal(t, model.TagTable, model.TagDB{}.TableName())
}
