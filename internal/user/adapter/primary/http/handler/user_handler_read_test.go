package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/platform/config"
	httperrors "github.com/lechitz/AionApi/internal/platform/server/http/errors"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/stretchr/testify/require"
)

func TestUserReadHandlers(t *testing.T) {
	t.Run("list all success", func(t *testing.T) {
		svc := &mockUserService{listAllFn: func(context.Context) ([]userdomain.User, error) {
			return []userdomain.User{{ID: 1, Username: "u1", Email: "u1@example.com", CreatedAt: time.Now().UTC()}}, nil
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/all", nil)
		rec := httptest.NewRecorder()

		h.ListAll(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "users retrieved successfully")
	})

	t.Run("list all service error", func(t *testing.T) {
		svc := &mockUserService{listAllFn: func(context.Context) ([]userdomain.User, error) {
			return nil, errors.New("boom")
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/all", nil)
		rec := httptest.NewRecorder()

		h.ListAll(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("get me missing user id", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/me", nil)
		rec := httptest.NewRecorder()

		h.GetMe(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("get me success", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/me", nil)
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()

		h.GetMe(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "user_me_success")
	})

	t.Run("get user by id missing param", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/anything", nil)
		rec := httptest.NewRecorder()
		h.GetUserByID(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("get user by id invalid param", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/abc", nil)
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Get("/user/{user_id}", h.GetUserByID)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("get user by id not found", func(t *testing.T) {
		svc := &mockUserService{getByIDFn: func(context.Context, uint64) (userdomain.User, error) {
			return userdomain.User{}, httperrors.ErrResourceNotFound
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/9", nil)
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Get("/user/{user_id}", h.GetUserByID)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("get user by id success", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/user/9", nil)
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Get("/user/{user_id}", h.GetUserByID)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "user retrieved successfully")
	})
}
