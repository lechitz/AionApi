package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	userinput "github.com/lechitz/AionApi/internal/user/core/ports/input"
	"github.com/stretchr/testify/require"
)

func TestCreateUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequest(
			http.MethodPost,
			"/user/create",
			strings.NewReader(`{"name":"John Doe","username":"john","email":"john@example.com","password":"12345678"}`),
		)
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		parsed := decodeResponseBody(t, rec.Body.Bytes())
		require.EqualValues(t, http.StatusCreated, parsed["code"])
		require.Equal(t, "user created successfully", parsed["message"])
	})

	t.Run("decode error", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(`{"name":`))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(`{"name":"John","username":"john","email":"john@example.com","password":"123"}`))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserService{createFn: func(context.Context, userinput.CreateUserCommand) (userdomain.User, error) {
			return userdomain.User{}, sharederrors.ErrDomainConflict
		}}
		h := handler.New(svc, &config.Config{}, mockLogger{})
		req := httptest.NewRequest(
			http.MethodPost,
			"/user/create",
			strings.NewReader(`{"name":"John Doe","username":"john","email":"john@example.com","password":"12345678"}`),
		)
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("large body", func(t *testing.T) {
		h := handler.New(&mockUserService{}, &config.Config{}, mockLogger{})
		req := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(strings.Repeat("a", (1<<20)+8)))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
