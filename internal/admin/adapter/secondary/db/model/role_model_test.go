package model_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/admin/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestRoleModelTableNames(t *testing.T) {
	require.Equal(t, model.TableRoles, model.RoleDB{}.TableName())
	require.Equal(t, model.TableUserRoles, model.UserRoleDB{}.TableName())
}
