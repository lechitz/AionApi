package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	userinput "github.com/lechitz/AionApi/internal/user/core/ports/input"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserHandler(t *testing.T) {
	t.Run("missing context", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user", strings.NewReader(`{"name":"New"}`))
		rec := httptest.NewRecorder()
		h.UpdateUser(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user", strings.NewReader(`{"name":`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUser(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("no fields", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user", strings.NewReader(`{}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUser(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserService{updateUserFn: func(context.Context, uint64, userinput.UpdateUserCommand) (userdomain.User, error) {
			return userdomain.User{}, errors.New("boom")
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user", strings.NewReader(`{"name":"New"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUser(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user", strings.NewReader(`{"name":"New Name"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUser(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "user updated successfully")
	})
}

func TestUpdatePasswordHandler(t *testing.T) {
	t.Run("missing required fields", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user/password", strings.NewReader(`{"password":"old"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUserPassword(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing context", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user/password", strings.NewReader(`{"password":"old","new_password":"new"}`))
		rec := httptest.NewRecorder()
		h.UpdateUserPassword(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("domain conflict", func(t *testing.T) {
		svc := &mockUserService{updatePasswordFn: func(context.Context, uint64, string, string) (string, error) {
			return "", sharederrors.ErrDomainConflict
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user/password", strings.NewReader(`{"password":"old","new_password":"new"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUserPassword(rec, req)
		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("success sets cookie", func(t *testing.T) {
		cfg := &config.Config{Cookie: config.CookieConfig{Path: "/", Domain: "localhost", SameSite: "Lax"}}
		h := handler.New(&mockUserService{}, cfg, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/user/password", strings.NewReader(`{"password":"old","new_password":"new"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.UpdateUserPassword(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.NotEmpty(t, rec.Header().Values("Set-Cookie"))
		require.Contains(t, rec.Header().Get("Set-Cookie"), commonkeys.AuthTokenCookieName)
	})
}

func TestSoftDeleteUserHandler(t *testing.T) {
	t.Run("missing context", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/user", nil)
		rec := httptest.NewRecorder()
		h.SoftDeleteUser(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserService{softDeleteUserFn: func(context.Context, uint64) error {
			return errors.New("boom")
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/user", nil)
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.SoftDeleteUser(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		cfg := &config.Config{Cookie: config.CookieConfig{Path: "/", Domain: "localhost", SameSite: "Lax"}}
		h := handler.New(&mockUserService{}, cfg, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/user", nil)
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()
		h.SoftDeleteUser(rec, req)
		require.Equal(t, http.StatusNoContent, rec.Code)
		require.NotEmpty(t, rec.Header().Values("Set-Cookie"))
	})
}
