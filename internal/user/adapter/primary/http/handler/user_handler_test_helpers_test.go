package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	authdomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	authinput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	userinput "github.com/lechitz/AionApi/internal/user/core/ports/input"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	createFn           func(ctx context.Context, cmd userinput.CreateUserCommand) (userdomain.User, error)
	getByIDFn          func(ctx context.Context, userID uint64) (userdomain.User, error)
	getByUsernameFn    func(ctx context.Context, username string) (userdomain.User, error)
	listAllFn          func(ctx context.Context) ([]userdomain.User, error)
	getUserStatsFn     func(ctx context.Context, userID uint64) (userdomain.UserStats, error)
	updateUserFn       func(ctx context.Context, userID uint64, cmd userinput.UpdateUserCommand) (userdomain.User, error)
	updatePasswordFn   func(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error)
	softDeleteUserFn   func(ctx context.Context, userID uint64) error
	calledUpdateUserID uint64
	calledPasswordID   uint64
}

func (m *mockUserService) Create(ctx context.Context, cmd userinput.CreateUserCommand) (userdomain.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, cmd)
	}
	return userdomain.User{ID: 1, Name: cmd.Name, Username: cmd.Username, Email: cmd.Email}, nil
}

func (m *mockUserService) GetByID(ctx context.Context, userID uint64) (userdomain.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, userID)
	}
	return userdomain.User{ID: userID, Name: "User", Username: "user", Email: "user@example.com"}, nil
}

func (m *mockUserService) GetUserByUsername(ctx context.Context, username string) (userdomain.User, error) {
	if m.getByUsernameFn != nil {
		return m.getByUsernameFn(ctx, username)
	}
	return userdomain.User{ID: 99, Username: username}, nil
}

func (m *mockUserService) ListAll(ctx context.Context) ([]userdomain.User, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx)
	}
	return []userdomain.User{}, nil
}

func (m *mockUserService) GetUserStats(ctx context.Context, userID uint64) (userdomain.UserStats, error) {
	if m.getUserStatsFn != nil {
		return m.getUserStatsFn(ctx, userID)
	}
	return userdomain.UserStats{}, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, userID uint64, cmd userinput.UpdateUserCommand) (userdomain.User, error) {
	m.calledUpdateUserID = userID
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, userID, cmd)
	}
	return userdomain.User{ID: userID, Username: "updated", Email: "updated@example.com", UpdatedAt: time.Now().UTC()}, nil
}

func (m *mockUserService) UpdatePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error) {
	m.calledPasswordID = userID
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, userID, oldPassword, newPassword)
	}
	return "new-token", nil
}

func (m *mockUserService) SoftDeleteUser(ctx context.Context, userID uint64) error {
	if m.softDeleteUserFn != nil {
		return m.softDeleteUserFn(ctx, userID)
	}
	return nil
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
	return 1, map[string]any{"roles": []string{"user"}}, nil
}

func (mockAuthService) Login(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
	return authdomain.AuthenticatedUser{}, "", "", nil
}

func (mockAuthService) Logout(context.Context, uint64) error { return nil }

func (mockAuthService) RefreshTokenRenewal(context.Context, string) (string, string, error) {
	return "", "", nil
}

type mockRouter struct {
	groups         []string
	posts          []string
	gets           []string
	puts           []string
	deletes        []string
	groupWithCalls int
}

func (m *mockRouter) Use(...ports.Middleware) {}
func (m *mockRouter) Group(prefix string, fn func(ports.Router)) {
	m.groups = append(m.groups, prefix)
	fn(m)
}

func (m *mockRouter) GroupWith(_ ports.Middleware, fn func(ports.Router)) {
	m.groupWithCalls++
	fn(m)
}
func (m *mockRouter) Mount(string, http.Handler)                               {}
func (m *mockRouter) Handle(string, string, http.Handler)                      {}
func (m *mockRouter) GET(path string, _ http.Handler)                          { m.gets = append(m.gets, path) }
func (m *mockRouter) POST(path string, _ http.Handler)                         { m.posts = append(m.posts, path) }
func (m *mockRouter) PUT(path string, _ http.Handler)                          { m.puts = append(m.puts, path) }
func (m *mockRouter) DELETE(path string, _ http.Handler)                       { m.deletes = append(m.deletes, path) }
func (m *mockRouter) SetNotFound(http.Handler)                                 {}
func (m *mockRouter) SetMethodNotAllowed(http.Handler)                         {}
func (m *mockRouter) SetError(func(http.ResponseWriter, *http.Request, error)) {}
func (m *mockRouter) ServeHTTP(http.ResponseWriter, *http.Request)             {}

func decodeResponseBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var parsed map[string]any
	require.NoError(t, json.Unmarshal(body, &parsed))
	return parsed
}

var (
	_ userinput.UserService = (*mockUserService)(nil)
	_ authinput.AuthService = (*mockAuthService)(nil)
)
