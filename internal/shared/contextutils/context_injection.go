// Package contextutils provides utilities for working with context.
package contextutils

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"

	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// InjectUserIntoContext puts userID from claims into context if present.
func InjectUserIntoContext(ctx context.Context, claims map[string]any) context.Context {
	userIDFloat, ok := claims[claimskeys.UserID].(float64)
	if !ok {
		return ctx
	}

	userID := uint64(userIDFloat)

	return context.WithValue(ctx, ctxkeys.UserID, userID)
}

// GetRequestID returns the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(ctxkeys.RequestID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetTraceID returns the trace ID from the context.
func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(ctxkeys.TraceID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetUserID returns the user ID from the context.
func GetUserID(ctx context.Context) string {
	if v := ctx.Value(ctxkeys.UserID); v != nil {
		switch id := v.(type) {
		case uint64:
			return strconv.FormatUint(id, 10)
		case int64:
			return strconv.FormatInt(id, 10)
		case string:
			return id
		}
	}
	return ""
}
