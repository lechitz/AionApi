package handler_test

import (
	"context"
	"net/http"
	"testing"

	handlerpkg "github.com/lechitz/aion-api/internal/audit/adapter/primary/http/handler"
	authdomain "github.com/lechitz/aion-api/internal/auth/core/domain"
	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
	"github.com/stretchr/testify/require"
)

type mockAuditRouter struct {
	groupWithCall int
	gets          []string
}

func (m *mockAuditRouter) Use(...ports.Middleware)          {}
func (m *mockAuditRouter) Group(string, func(ports.Router)) {}
func (m *mockAuditRouter) GroupWith(_ ports.Middleware, fn func(ports.Router)) {
	m.groupWithCall++
	fn(m)
}
func (m *mockAuditRouter) Mount(string, http.Handler)                               {}
func (m *mockAuditRouter) Handle(string, string, http.Handler)                      {}
func (m *mockAuditRouter) GET(path string, _ http.Handler)                          { m.gets = append(m.gets, path) }
func (m *mockAuditRouter) POST(string, http.Handler)                                {}
func (m *mockAuditRouter) PUT(string, http.Handler)                                 {}
func (m *mockAuditRouter) DELETE(string, http.Handler)                              {}
func (m *mockAuditRouter) SetNotFound(http.Handler)                                 {}
func (m *mockAuditRouter) SetMethodNotAllowed(http.Handler)                         {}
func (m *mockAuditRouter) SetError(func(http.ResponseWriter, *http.Request, error)) {}
func (m *mockAuditRouter) ServeHTTP(http.ResponseWriter, *http.Request)             {}

type authServiceStub struct{}

func (authServiceStub) Login(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
	return authdomain.AuthenticatedUser{}, "", "", nil
}

func (authServiceStub) Validate(context.Context, string) (uint64, map[string]any, error) {
	return 0, nil, nil
}

func (authServiceStub) Logout(context.Context, uint64) error { return nil }

func (authServiceStub) RefreshTokenRenewal(context.Context, string) (string, string, error) {
	return "", "", nil
}

func TestRegisterHTTP(t *testing.T) {
	h, _ := newAuditHandler(t)
	router := &mockAuditRouter{}

	handlerpkg.RegisterHTTP(router, h, authServiceStub{}, nil)

	require.Equal(t, 1, router.groupWithCall)
	require.Equal(t, []string{"/audit/events"}, router.gets)
}

func TestRegisterHTTP_NoAuthService(t *testing.T) {
	h, _ := newAuditHandler(t)
	router := &mockAuditRouter{}

	handlerpkg.RegisterHTTP(router, h, nil, nil)

	require.Equal(t, 0, router.groupWithCall)
	require.Empty(t, router.gets)
}
