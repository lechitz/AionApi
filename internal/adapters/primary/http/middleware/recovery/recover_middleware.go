// Package recovery provides HTTP middleware that recovers from panics,
// logs the error with stack trace, and responds with a 500 error.
package recovery

import (
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
)

// RecoverMiddleware is a middleware that recovers from panics, logs the error, and returns an internal server error response.
func RecoverMiddleware(log output.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					errorID := uuid.New().String()

					// TODO: ajustar magic strings.
					log.Errorw("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"method", r.Method,
						"stack", string(debug.Stack()),
						"error_id", errorID,
					)

					response.HandleError(w, log, http.StatusInternalServerError,
						"unexpected server error (ref: "+errorID+")", nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
