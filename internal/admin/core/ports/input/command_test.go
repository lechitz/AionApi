package input_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/admin/core/ports/input"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserRolesCommandHasUpdates(t *testing.T) {
	require.False(t, input.UpdateUserRolesCommand{}.HasUpdates())
	require.True(t, input.UpdateUserRolesCommand{Roles: []string{"admin"}}.HasUpdates())
}
