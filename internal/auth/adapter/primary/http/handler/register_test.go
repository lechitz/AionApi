package handler_test

import (
	"net/http"
	"testing"

	handlerpkg "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	"github.com/stretchr/testify/require"
)

type mockAuthRouter struct {
	groupPrefixes []string
	groupWithCall int
	gets          []string
	posts         []string
}

func (m *mockAuthRouter) Use(...ports.Middleware) {}
func (m *mockAuthRouter) Group(prefix string, fn func(ports.Router)) {
	m.groupPrefixes = append(m.groupPrefixes, prefix)
	fn(m)
}

func (m *mockAuthRouter) GroupWith(_ ports.Middleware, fn func(ports.Router)) {
	m.groupWithCall++
	fn(m)
}
func (m *mockAuthRouter) Mount(string, http.Handler)                               {}
func (m *mockAuthRouter) Handle(string, string, http.Handler)                      {}
func (m *mockAuthRouter) GET(path string, _ http.Handler)                          { m.gets = append(m.gets, path) }
func (m *mockAuthRouter) POST(path string, _ http.Handler)                         { m.posts = append(m.posts, path) }
func (m *mockAuthRouter) PUT(string, http.Handler)                                 {}
func (m *mockAuthRouter) DELETE(string, http.Handler)                              {}
func (m *mockAuthRouter) SetNotFound(http.Handler)                                 {}
func (m *mockAuthRouter) SetMethodNotAllowed(http.Handler)                         {}
func (m *mockAuthRouter) SetError(func(http.ResponseWriter, *http.Request, error)) {}
func (m *mockAuthRouter) ServeHTTP(http.ResponseWriter, *http.Request)             {}

func TestRegisterHTTP(t *testing.T) {
	h := newTestHandler(t, &authServiceStub{})
	router := &mockAuthRouter{}

	handlerpkg.RegisterHTTP(router, h)

	require.Equal(t, []string{"/auth"}, router.groupPrefixes)
	require.Equal(t, 1, router.groupWithCall)
	require.ElementsMatch(t, []string{"/login", "/refresh", "/logout"}, router.posts)
	require.ElementsMatch(t, []string{"/session"}, router.gets)
}
