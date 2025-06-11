// Package recovery provides HTTP middleware that recovers from panics,
// logs the error with stack trace, and responds with a 500 error.
package recovery

import (
	"net/http"
	"runtime/debug"

	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// RecoverMiddleware is a middleware that recovers from panics, logs the error, and returns an internal server error response.
func RecoverMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Errorw("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"method", r.Method,
						"stack", string(debug.Stack()),
					)

					response.HandleError(
						w,
						log,
						http.StatusInternalServerError,
						"internal server error",
						nil,
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
