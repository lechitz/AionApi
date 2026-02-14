package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authdomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	authinput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	handler "github.com/lechitz/AionApi/internal/chat/adapter/primary/http/handler"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	chatinput "github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
)

type mockChatService struct {
	processFn func(ctx context.Context, userID uint64, message string, requestContext map[string]interface{}) (*domain.ChatResult, error)
}

func (m mockChatService) ProcessMessage(ctx context.Context, userID uint64, message string, requestContext map[string]interface{}) (*domain.ChatResult, error) {
	if m.processFn != nil {
		return m.processFn(ctx, userID, message, requestContext)
	}
	return &domain.ChatResult{}, nil
}

func (mockChatService) SaveChatHistory(context.Context, uint64, string, string, int, map[string]string) error {
	return nil
}

func (mockChatService) GetChatHistory(context.Context, uint64, int, int) ([]domain.ChatHistory, error) {
	return nil, nil
}

func (mockChatService) GetLatestChatHistory(context.Context, uint64, int) ([]domain.ChatHistory, error) {
	return nil, nil
}

func (mockChatService) GetChatContext(context.Context, uint64) (*domain.ChatContext, error) {
	return &domain.ChatContext{}, nil
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
	return 0, nil, nil
}

func (mockAuthService) Login(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
	return authdomain.AuthenticatedUser{}, "", "", nil
}
func (mockAuthService) Logout(context.Context, uint64) error { return nil }
func (mockAuthService) RefreshTokenRenewal(context.Context, string) (string, string, error) {
	return "", "", nil
}

type mockRouter struct {
	groupWithCalls int
	posts          []string
}

func (m *mockRouter) Use(...ports.Middleware)          {}
func (m *mockRouter) Group(string, func(ports.Router)) {}
func (m *mockRouter) GroupWith(_ ports.Middleware, fn func(ports.Router)) {
	m.groupWithCalls++
	fn(m)
}
func (m *mockRouter) Mount(string, http.Handler)                               {}
func (m *mockRouter) Handle(string, string, http.Handler)                      {}
func (m *mockRouter) GET(string, http.Handler)                                 {}
func (m *mockRouter) POST(path string, _ http.Handler)                         { m.posts = append(m.posts, path) }
func (m *mockRouter) PUT(string, http.Handler)                                 {}
func (m *mockRouter) DELETE(string, http.Handler)                              {}
func (m *mockRouter) SetNotFound(http.Handler)                                 {}
func (m *mockRouter) SetMethodNotAllowed(http.Handler)                         {}
func (m *mockRouter) SetError(func(http.ResponseWriter, *http.Request, error)) {}
func (m *mockRouter) ServeHTTP(http.ResponseWriter, *http.Request)             {}

func TestNewAndRegisterHTTP(t *testing.T) {
	h := handler.New(mockChatService{}, &config.Config{}, mockLogger{})
	require.NotNil(t, h)

	router := &mockRouter{}
	handler.RegisterHTTP(router, h, nil, mockLogger{})
	require.Equal(t, 0, router.groupWithCalls)
	require.Empty(t, router.posts)

	handler.RegisterHTTP(router, h, mockAuthService{}, mockLogger{})
	require.Equal(t, 1, router.groupWithCalls)
	require.ElementsMatch(t, []string{"/chat/text", "/chat/audio"}, router.posts)
}

func TestChatText_Success(t *testing.T) {
	h := handler.New(mockChatService{
		processFn: func(_ context.Context, userID uint64, message string, requestContext map[string]interface{}) (*domain.ChatResult, error) {
			require.Equal(t, uint64(7), userID)
			require.Equal(t, "hello", message)
			require.Equal(t, "v", requestContext["k"])
			return &domain.ChatResult{
				Response:   "ok",
				UI:         map[string]interface{}{"type": "simple"},
				Sources:    []interface{}{map[string]interface{}{"id": 1}, "ignored"},
				TokensUsed: 12,
			}, nil
		},
	}, &config.Config{}, mockLogger{})

	req := httptest.NewRequest(http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello","context":{"k":"v"}}`))
	req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ChatText(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Chat processed successfully")
	require.Contains(t, rec.Body.String(), "\"total_tokens\":12")
	require.Contains(t, rec.Body.String(), "\"sources\":[{")
}

func TestChatText_Errors(t *testing.T) {
	h := handler.New(mockChatService{processFn: func(context.Context, uint64, string, map[string]interface{}) (*domain.ChatResult, error) {
		return nil, errors.New("boom")
	}}, &config.Config{}, mockLogger{})

	t.Run("missing user id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid user id type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, "7"))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat/text", strings.NewReader(`{"message":`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat/text", strings.NewReader(`{"message":"   "}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

var (
	_ chatinput.ChatService = (*mockChatService)(nil)
	_ authinput.AuthService = (*mockAuthService)(nil)
)
