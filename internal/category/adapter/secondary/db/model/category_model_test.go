package model_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestCategoryModelTableName(t *testing.T) {
	require.Equal(t, model.CategoryTable, model.CategoryDB{}.TableName())
}
