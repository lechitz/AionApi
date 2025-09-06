// Package middleware are used to define custom middleware for the GraphQL schema
package middleware

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

func Auth() func(ctx context.Context, obj interface{}, next graphql.Resolver, role *string) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, role *string) (interface{}, error) {
		userID, _ := ctx.Value(ctxkeys.UserID).(uint64)
		if userID == 0 {
			return nil, errors.New("unauthenticated")
		}
		// Opcional: checar papel quando ctx trouxer roles do middleware
		// if role != nil && !hasRole(ctx, *role) { return nil, errors.New("forbidden") }
		return next(ctx)
	}
}
