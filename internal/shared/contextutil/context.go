// Package contextutil provides helper functions for working with context.Context values.
package contextutil

import (
	"context"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	TokenKey  contextKey = "token"
)

// GetUserID retrieves the user ID from the context, if available.
func GetUserID(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint64)
	return userID, ok
}

// GetToken retrieves the authentication token from the context, if available.
func GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(TokenKey).(string)
	return token, ok
}
