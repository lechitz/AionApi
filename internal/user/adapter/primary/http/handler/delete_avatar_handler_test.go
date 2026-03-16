package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/stretchr/testify/require"
)

func TestDeleteAvatarHandler(t *testing.T) {
	t.Run("missing context", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/user/avatar", nil)
		rec := httptest.NewRecorder()

		h.DeleteAvatar(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserService{removeAvatarFn: func(context.Context, uint64) (userdomain.User, error) {
			return userdomain.User{}, errors.New("boom")
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/user/avatar", nil)
		req = req.WithContext(context.WithValue(req.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()

		h.DeleteAvatar(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		svc := &mockUserService{removeAvatarFn: func(_ context.Context, userID uint64) (userdomain.User, error) {
			return userdomain.User{
				ID:                  userID,
				Name:                "Demo User",
				Username:            "testuser",
				Email:               "test@aion.local",
				AvatarURL:           nil,
				OnboardingCompleted: true,
				UpdatedAt:           time.Now().UTC(),
			}, nil
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/user/avatar", nil)
		req = req.WithContext(context.WithValue(req.Context(), ctxkeys.UserID, uint64(10)))
		rec := httptest.NewRecorder()

		h.DeleteAvatar(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "user avatar deleted successfully")
		require.NotContains(t, rec.Body.String(), "\"avatar_url\":\"")
	})
}
