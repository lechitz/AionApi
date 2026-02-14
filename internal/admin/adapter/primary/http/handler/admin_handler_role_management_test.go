package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	handler "github.com/lechitz/AionApi/internal/admin/adapter/primary/http/handler"
	admin "github.com/lechitz/AionApi/internal/admin/core/domain"
	admininput "github.com/lechitz/AionApi/internal/admin/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
)

func makeRoleReq(t *testing.T, url string, withActor bool, claims map[string]any) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, url, nil)
	ctx := t.Context()
	if withActor {
		ctx = context.WithValue(ctx, ctxkeys.UserID, uint64(100))
	}
	if claims != nil {
		ctx = context.WithValue(ctx, ctxkeys.Claims, claims)
	}
	return req.WithContext(ctx)
}

func serveRoleRoute(path string, fn http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	r.Put(path, fn)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestPromoteToAdmin(t *testing.T) {
	h := handler.New(mockAdminService{}, &config.Config{}, mockLogger{})
	rec := serveRoleRoute("/admin/users/{user_id}/promote-admin", h.PromoteToAdmin, makeRoleReq(t, "/admin/users/5/promote-admin", false, nil))
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = serveRoleRoute("/admin/users/{user_id}/promote-admin", h.PromoteToAdmin, makeRoleReq(t, "/admin/users/abc/promote-admin", true, nil))
	require.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serveRoleRoute(
		"/admin/users/{user_id}/promote-admin",
		h.PromoteToAdmin,
		makeRoleReq(t, "/admin/users/5/promote-admin", true, map[string]any{commonkeys.Roles: []any{"owner"}}),
	)
	require.Equal(t, http.StatusOK, rec.Code)

	h = handler.New(mockAdminService{promoteFn: func(context.Context, admininput.PromoteToAdminCommand) (admin.AdminUser, error) {
		return admin.AdminUser{}, errors.New("boom")
	}}, &config.Config{}, mockLogger{})
	rec = serveRoleRoute(
		"/admin/users/{user_id}/promote-admin",
		h.PromoteToAdmin,
		makeRoleReq(t, "/admin/users/5/promote-admin", true, map[string]any{commonkeys.Roles: []string{"admin"}}),
	)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDemoteFromAdmin(t *testing.T) {
	h := handler.New(mockAdminService{}, &config.Config{}, mockLogger{})
	rec := serveRoleRoute("/admin/users/{user_id}/demote-admin", h.DemoteFromAdmin, makeRoleReq(t, "/admin/users/5/demote-admin", false, nil))
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = serveRoleRoute("/admin/users/{user_id}/demote-admin", h.DemoteFromAdmin, makeRoleReq(t, "/admin/users/abc/demote-admin", true, nil))
	require.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serveRoleRoute(
		"/admin/users/{user_id}/demote-admin",
		h.DemoteFromAdmin,
		makeRoleReq(t, "/admin/users/5/demote-admin", true, map[string]any{commonkeys.Roles: []string{"admin"}}),
	)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestBlockAndUnblock(t *testing.T) {
	h := handler.New(mockAdminService{}, &config.Config{}, mockLogger{})
	rec := serveRoleRoute("/admin/users/{user_id}/block", h.BlockUser, makeRoleReq(t, "/admin/users/5/block", false, nil))
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = serveRoleRoute("/admin/users/{user_id}/block", h.BlockUser, makeRoleReq(t, "/admin/users/5/block", true, map[string]any{commonkeys.Roles: []string{"admin"}}))
	require.Equal(t, http.StatusOK, rec.Code)

	h = handler.New(mockAdminService{blockFn: func(context.Context, admininput.BlockUserCommand) (admin.AdminUser, error) {
		return admin.AdminUser{}, errors.New("block-fail")
	}}, &config.Config{}, mockLogger{})
	rec = serveRoleRoute("/admin/users/{user_id}/block", h.BlockUser, makeRoleReq(t, "/admin/users/5/block", true, map[string]any{commonkeys.Roles: []any{"admin"}}))
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	h = handler.New(mockAdminService{}, &config.Config{}, mockLogger{})
	rec = serveRoleRoute(
		"/admin/users/{user_id}/unblock",
		h.UnblockUser,
		makeRoleReq(t, "/admin/users/5/unblock", true, map[string]any{commonkeys.Roles: []any{"admin"}}),
	)
	require.Equal(t, http.StatusOK, rec.Code)

	h = handler.New(mockAdminService{unblockFn: func(context.Context, admininput.UnblockUserCommand) (admin.AdminUser, error) {
		return admin.AdminUser{}, errors.New("unblock-fail")
	}}, &config.Config{}, mockLogger{})
	rec = serveRoleRoute(
		"/admin/users/{user_id}/unblock",
		h.UnblockUser,
		makeRoleReq(t, "/admin/users/5/unblock", true, map[string]any{commonkeys.Roles: []string{"admin"}}),
	)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
