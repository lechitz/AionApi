package model_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestRecordModelTableName(t *testing.T) {
	require.Equal(t, "aion_api.records", model.Record{}.TableName())
}
