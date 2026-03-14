package model_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestUserModelTableName(t *testing.T) {
	require.Equal(t, model.TableUsers, model.UserDB{}.TableName())
}
