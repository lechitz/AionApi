package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	admin "github.com/lechitz/AionApi/internal/admin/core/domain"
	admininput "github.com/lechitz/AionApi/internal/admin/core/ports/input"
	authdomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	authinput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	"github.com/stretchr/testify/require"
)

type mockAdminService struct {
	updateUserRolesFn func(context.Context, admininput.UpdateUserRolesCommand) (admin.AdminUser, error)
	promoteFn         func(context.Context, admininput.PromoteToAdminCommand) (admin.AdminUser, error)
	demoteFn          func(context.Context, admininput.DemoteFromAdminCommand) (admin.AdminUser, error)
	blockFn           func(context.Context, admininput.BlockUserCommand) (admin.AdminUser, error)
	unblockFn         func(context.Context, admininput.UnblockUserCommand) (admin.AdminUser, error)
}

func (m mockAdminService) UpdateUserRoles(ctx context.Context, cmd admininput.UpdateUserRolesCommand) (admin.AdminUser, error) {
	if m.updateUserRolesFn != nil {
		return m.updateUserRolesFn(ctx, cmd)
	}
	return admin.AdminUser{ID: cmd.UserID, Username: "u", Email: "u@example.com", Roles: cmd.Roles, UpdatedAt: time.Now().UTC()}, nil
}

func (m mockAdminService) PromoteToAdmin(ctx context.Context, cmd admininput.PromoteToAdminCommand) (admin.AdminUser, error) {
	if m.promoteFn != nil {
		return m.promoteFn(ctx, cmd)
	}
	return admin.AdminUser{ID: cmd.UserID, Username: "u", Email: "u@example.com", Roles: []string{"admin"}, UpdatedAt: time.Now().UTC()}, nil
}

func (m mockAdminService) DemoteFromAdmin(ctx context.Context, cmd admininput.DemoteFromAdminCommand) (admin.AdminUser, error) {
	if m.demoteFn != nil {
		return m.demoteFn(ctx, cmd)
	}
	return admin.AdminUser{ID: cmd.UserID, Username: "u", Email: "u@example.com", Roles: []string{"user"}, UpdatedAt: time.Now().UTC()}, nil
}

func (m mockAdminService) BlockUser(ctx context.Context, cmd admininput.BlockUserCommand) (admin.AdminUser, error) {
	if m.blockFn != nil {
		return m.blockFn(ctx, cmd)
	}
	return admin.AdminUser{ID: cmd.UserID, Username: "u", Email: "u@example.com", Roles: []string{"blocked"}, UpdatedAt: time.Now().UTC()}, nil
}

func (m mockAdminService) UnblockUser(ctx context.Context, cmd admininput.UnblockUserCommand) (admin.AdminUser, error) {
	if m.unblockFn != nil {
		return m.unblockFn(ctx, cmd)
	}
	return admin.AdminUser{ID: cmd.UserID, Username: "u", Email: "u@example.com", Roles: []string{"user"}, UpdatedAt: time.Now().UTC()}, nil
}

type mockLogger struct{}

func (mockLogger) Infof(string, ...any)                      {}
func (mockLogger) Errorf(string, ...any)                     {}
func (mockLogger) Debugf(string, ...any)                     {}
func (mockLogger) Warnf(string, ...any)                      {}
func (mockLogger) Infow(string, ...any)                      {}
func (mockLogger) Errorw(string, ...any)                     {}
func (mockLogger) Debugw(string, ...any)                     {}
func (mockLogger) Warnw(string, ...any)                      {}
func (mockLogger) InfowCtx(context.Context, string, ...any)  {}
func (mockLogger) ErrorwCtx(context.Context, string, ...any) {}
func (mockLogger) WarnwCtx(context.Context, string, ...any)  {}
func (mockLogger) DebugwCtx(context.Context, string, ...any) {}

type mockAuthService struct{}

func (mockAuthService) Validate(context.Context, string) (uint64, map[string]any, error) {
	return 1, map[string]any{"roles": []string{"admin"}}, nil
}

func (mockAuthService) Login(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
	return authdomain.AuthenticatedUser{}, "", "", nil
}
func (mockAuthService) Logout(context.Context, uint64) error { return nil }
func (mockAuthService) RefreshTokenRenewal(context.Context, string) (string, string, error) {
	return "", "", nil
}

type mockRouter struct {
	puts           []string
	groupWithCalls int
}

func (m *mockRouter) Use(...ports.Middleware) {}
func (m *mockRouter) Group(string, func(ports.Router)) {
}

func (m *mockRouter) GroupWith(_ ports.Middleware, fn func(ports.Router)) {
	m.groupWithCalls++
	fn(m)
}
func (m *mockRouter) Mount(string, http.Handler)                               {}
func (m *mockRouter) Handle(string, string, http.Handler)                      {}
func (m *mockRouter) GET(string, http.Handler)                                 {}
func (m *mockRouter) POST(string, http.Handler)                                {}
func (m *mockRouter) PUT(path string, _ http.Handler)                          { m.puts = append(m.puts, path) }
func (m *mockRouter) DELETE(string, http.Handler)                              {}
func (m *mockRouter) SetNotFound(http.Handler)                                 {}
func (m *mockRouter) SetMethodNotAllowed(http.Handler)                         {}
func (m *mockRouter) SetError(func(http.ResponseWriter, *http.Request, error)) {}
func (m *mockRouter) ServeHTTP(http.ResponseWriter, *http.Request)             {}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var parsed map[string]any
	require.NoError(t, json.Unmarshal(body, &parsed))
	return parsed
}

var (
	_ admininput.AdminService = (*mockAdminService)(nil)
	_ authinput.AuthService   = (*mockAuthService)(nil)
)
