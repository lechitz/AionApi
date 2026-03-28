package directives_test

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/directives"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
)

func TestAuthDirective(t *testing.T) {
	dir := directives.Auth()

	next := func(_ context.Context) (interface{}, error) {
		return "ok", nil
	}

	t.Run("missing user id", func(t *testing.T) {
		res, err := dir(t.Context(), nil, graphql.Resolver(next), nil)
		require.Nil(t, res)
		require.Error(t, err)
		require.Equal(t, 401, sharederrors.MapErrorToHTTPStatus(err))
	})

	t.Run("service account bypasses role check", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
		ctx = context.WithValue(ctx, ctxkeys.ServiceAccount, true)
		roles := "admin"

		res, err := dir(ctx, nil, graphql.Resolver(next), &roles)
		require.NoError(t, err)
		require.Equal(t, "ok", res)
	})

	t.Run("roles as []string", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
		ctx = context.WithValue(ctx, ctxkeys.Claims, map[string]any{commonkeys.Roles: []string{"user", "admin"}})
		roles := "admin"

		res, err := dir(ctx, nil, graphql.Resolver(next), &roles)
		require.NoError(t, err)
		require.Equal(t, "ok", res)
	})

	t.Run("roles as []any", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
		ctx = context.WithValue(ctx, ctxkeys.Claims, []any{"user", "admin"})
		roles := "admin"

		res, err := dir(ctx, nil, graphql.Resolver(next), &roles)
		require.NoError(t, err)
		require.Equal(t, "ok", res)
	})

	t.Run("roles as csv string", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
		ctx = context.WithValue(ctx, ctxkeys.Claims, "user, admin")
		roles := "admin"

		res, err := dir(ctx, nil, graphql.Resolver(next), &roles)
		require.NoError(t, err)
		require.Equal(t, "ok", res)
	})

	t.Run("forbidden when role missing", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
		ctx = context.WithValue(ctx, ctxkeys.Claims, map[string]any{commonkeys.Roles: []string{"user"}})
		roles := "admin"

		res, err := dir(ctx, nil, graphql.Resolver(next), &roles)
		require.Nil(t, res)
		require.Error(t, err)
		require.Equal(t, 403, sharederrors.MapErrorToHTTPStatus(err))
	})

	t.Run("no role requirement", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
		empty := ""
		res, err := dir(ctx, nil, graphql.Resolver(next), &empty)
		require.NoError(t, err)
		require.Equal(t, "ok", res)
	})
}
