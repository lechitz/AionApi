// Package directives contains custom directives for gqlgen.
package directives

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// Auth implements the @auth(role: String) directive.
// It enforces that a user is present in the context (populated by the HTTP auth middleware)
// and, optionally, that the user has the required "role".
func Auth() func(ctx context.Context, obj any, next graphql.Resolver, roles *string) (res any, err error) {
	return func(ctx context.Context, _ any, next graphql.Resolver, roles *string) (any, error) {
		// 1) Must be authenticated (user_id in context)
		if ctx.Value(ctxkeys.UserID) == nil {
			return nil, sharederrors.ErrUnauthorized("missing user in context")
		}

		// 2) (Optional) Check for a required role
		if roles != nil && *roles != "" {
			if !hasRole(ctx, *roles) {
				return nil, sharederrors.ErrForbidden("missing required roles: " + *roles)
			}
		}

		return next(ctx)
	}
}

// hasRole checks whether the context claims contain the required role.
// It supports common formats, e.g. claims["roles"] as []string, []any, or a CSV string.
func hasRole(ctx context.Context, required string) bool {
	rolesVal := extractRolesFromClaims(ctx.Value(ctxkeys.Claims))
	return rolesContain(rolesVal, required)
}

// extractRolesFromClaims reads "roles" from claims in common shapes.
// If claims is already a roles-like value (e.g. []string or string CSV), it returns it as-is.
func extractRolesFromClaims(claims any) any {
	switch m := claims.(type) {
	case map[string]any:
		return m["roles"]
	case []string:
		return m
	case []any:
		return m
	case string:
		return m
	default:
		return nil
	}
}

// rolesContain parses the roles value and checks if it contains the required role.
func rolesContain(v any, required string) bool {
	switch vv := v.(type) {
	case []string:
		for _, r := range vv {
			if r == required {
				return true
			}
		}
	case []any:
		for _, r := range vv {
			if s, ok := r.(string); ok && s == required {
				return true
			}
		}
	case string:
		// Support CSV in a single string: "admin,user,editor"
		for _, r := range strings.Split(vv, ",") {
			if strings.TrimSpace(r) == required {
				return true
			}
		}
	}
	return false
}
