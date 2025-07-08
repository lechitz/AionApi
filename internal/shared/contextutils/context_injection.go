// Package contextutils provides utilities for working with context.
package contextutils

import (
	"context"

	"github.com/lechitz/AionApi/internal/shared/ctxkeys"
)

// InjectUserIntoContext puts userID from claims into context if present.
func InjectUserIntoContext(ctx context.Context, claims map[string]any) context.Context {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return ctx
	}
	userID := uint64(userIDFloat)
	return context.WithValue(ctx, ctxkeys.UserID, userID)
}
