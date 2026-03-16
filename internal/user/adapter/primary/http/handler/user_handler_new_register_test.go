package handler_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	"github.com/stretchr/testify/require"
)

func TestNewAndRegisterHTTP(t *testing.T) {
	svc := &mockUserService{}
	h := handler.New(svc, &config.Config{}, mockLogger{})
	require.NotNil(t, h)

	r := &mockRouter{}
	handler.RegisterHTTP(r, h, nil, mockLogger{})
	require.Equal(t, []string{"/user", "/registration"}, r.groups)
	require.Equal(t, []string{"/create", "/avatar/upload", "/start", "/{registration_id}/complete"}, r.posts)
	require.Equal(t, []string{"/{registration_id}/profile", "/{registration_id}/avatar"}, r.puts)
	require.Equal(t, 0, r.groupWithCalls)

	r = &mockRouter{}
	handler.RegisterHTTP(r, h, mockAuthService{}, mockLogger{})
	require.Equal(t, []string{"/user", "/registration"}, r.groups)
	require.Equal(t, []string{"/create", "/avatar/upload", "/start", "/{registration_id}/complete"}, r.posts)
	require.Equal(t, []string{"/all", "/me", "/{user_id}"}, r.gets)
	require.Equal(t, []string{"/", "/password", "/{registration_id}/profile", "/{registration_id}/avatar"}, r.puts)
	require.Equal(t, []string{"/avatar", "/"}, r.deletes)
	require.Equal(t, 1, r.groupWithCalls)
}
