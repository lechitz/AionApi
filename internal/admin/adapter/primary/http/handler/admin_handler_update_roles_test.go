package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/aion-api/internal/admin/adapter/primary/http/dto"
	handler "github.com/lechitz/aion-api/internal/admin/adapter/primary/http/handler"
	admin "github.com/lechitz/aion-api/internal/admin/core/domain"
	admininput "github.com/lechitz/aion-api/internal/admin/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserRoles(t *testing.T) {
	newHandler := func(svc mockAdminService) *handler.Handler {
		return handler.New(svc, &config.Config{}, mockLogger{})
	}

	t.Run("missing user id", func(t *testing.T) {
		h := newHandler(mockAdminService{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users//roles", strings.NewReader(`{"roles":["admin"]}`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid user id", func(t *testing.T) {
		h := newHandler(mockAdminService{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users/abc/roles", strings.NewReader(`{"roles":["admin"]}`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		h := newHandler(mockAdminService{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users/1/roles", strings.NewReader(`{"roles":`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		h := newHandler(mockAdminService{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users/1/roles", strings.NewReader(`{"roles":[]}`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		h := newHandler(mockAdminService{updateUserRolesFn: func(context.Context, admininput.UpdateUserRolesCommand) (admin.AdminUser, error) {
			return admin.AdminUser{}, sharederrors.ErrDomainConflict
		}})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users/1/roles", strings.NewReader(`{"roles":["admin"]}`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		h := newHandler(mockAdminService{})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users/1/roles", strings.NewReader(`{"roles":["admin"]}`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		parsed := decodeBody(t, rec.Body.Bytes())
		require.EqualValues(t, http.StatusOK, parsed["code"])
		require.Equal(t, "user roles updated successfully", parsed["message"])
	})

	t.Run("dto validation sanity", func(t *testing.T) {
		err := dto.UpdateUserRolesRequest{Roles: []string{""}}.Validate()
		require.Error(t, err)
	})

	t.Run("unexpected service error", func(t *testing.T) {
		h := newHandler(mockAdminService{updateUserRolesFn: func(context.Context, admininput.UpdateUserRolesCommand) (admin.AdminUser, error) {
			return admin.AdminUser{}, errors.New("boom")
		}})
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/admin/users/1/roles", strings.NewReader(`{"roles":["admin"]}`))
		rec := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Put("/admin/users/{user_id}/roles", h.UpdateUserRoles)
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
