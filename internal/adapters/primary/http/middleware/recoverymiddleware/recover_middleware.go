// Package recoverymiddleware provides HTTP middleware that recovers from panics,
// logs the error with stack trace, and responds with a 500 error.
package recoverymiddleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"

	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
)

// New is a middleware that recovers from panics, logs the error, and returns an internal server error response.
func New(logger output.ContextLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					errorID := uuid.New().String()

					logger.Errorw("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"method", r.Method,
						"stack", string(debug.Stack()),
						"error_id", errorID,
					)

					panErr := errors.New("internal server error (ref: " + errorID + ")")

					httpresponse.WriteError(w, panErr, "unexpected server error", logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
