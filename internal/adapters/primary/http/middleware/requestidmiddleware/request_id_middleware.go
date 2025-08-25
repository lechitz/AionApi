// Package requestidmiddleware ensures a valid X-Request-ID header and context value.
package requestidmiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// New returns a middleware that injects/normalizes the X-Request-ID.
// - If no header is present, or it is invalid, a new UUID is generated.
// - If the header is present but longer than maxLen or not a UUID, it is replaced.
// - The final value is injected into both the request context and the response header.
func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := strings.TrimSpace(r.Header.Get(commonkeys.XRequestID))
			const maxLen = 128

			// Ensure a valid request ID: must be a UUID and not longer than maxLen.
			if reqID == "" || len(reqID) > maxLen || !isUUID(reqID) {
				reqID = uuid.NewString()
			} else if len(reqID) > maxLen {
				// If it's too long, truncate (should not happen with UUIDs but kept defensively).
				reqID = reqID[:maxLen]
			}

			// Inject into context and set a response header
			ctx := context.WithValue(r.Context(), ctxkeys.RequestID, reqID)
			w.Header().Set(commonkeys.XRequestID, reqID)

			// Continue request processing with enriched context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isUUID checks if a string is a valid UUID format.
func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
