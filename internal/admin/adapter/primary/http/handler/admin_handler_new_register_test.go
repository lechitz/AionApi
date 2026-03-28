package handler_test

import (
	"testing"

	handler "github.com/lechitz/aion-api/internal/admin/adapter/primary/http/handler"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func TestNewAndRegisterHTTP(t *testing.T) {
	h := handler.New(mockAdminService{}, &config.Config{}, mockLogger{})
	require.NotNil(t, h)

	r := &mockRouter{}
	handler.RegisterHTTP(r, h, nil, mockLogger{})
	require.Equal(t, 0, r.groupWithCalls)
	require.Empty(t, r.puts)

	r = &mockRouter{}
	handler.RegisterHTTP(r, h, mockAuthService{}, mockLogger{})
	require.Equal(t, 1, r.groupWithCalls)
	require.Equal(t, []string{
		"/admin/users/{user_id}/roles",
		"/admin/users/{user_id}/promote-admin",
		"/admin/users/{user_id}/demote-admin",
		"/admin/users/{user_id}/block",
		"/admin/users/{user_id}/unblock",
	}, r.puts)
}
