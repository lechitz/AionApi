// Package httpresponse centralizes HTTP response helpers for all handlers/domains.
package httpresponse

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
)

// ResponseBody defines the standard structure for all JSON API responses.
type ResponseBody struct {
	Date    time.Time `json:"date"`
	Result  any       `json:"result,omitempty"`
	Message string    `json:"message,omitempty"`
	Error   string    `json:"error,omitempty"`
	Details string    `json:"details,omitempty"`
	Code    int       `json:"code"`
}

// WriteJSON encodes any payload as JSON and writes it to the response with proper headers and status code.
// Allows passing extra headers (for CORS, Location, etc.).
func WriteJSON(w http.ResponseWriter, status int, payload any, headers ...map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	for _, hdr := range headers {
		for k, v := range hdr {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteSuccess sends a standardized success response with an optional result and message.
// Allows extra headers for advanced use cases.
func WriteSuccess(w http.ResponseWriter, status int, result any, message string, headers ...map[string]string) {
	body := ResponseBody{
		Date:    time.Now().UTC(),
		Result:  result,
		Message: message,
		Code:    status,
	}
	WriteJSON(w, status, body, headers...)
}

// WriteError sends a standardized error response, maps the error to HTTP status, and logs the error if logger is set.
func WriteError(w http.ResponseWriter, err error, message string, logger output.ContextLogger, headers ...map[string]string) {
	status := sharederrors.MapErrorToHTTPStatus(err)
	body := ResponseBody{
		Date:    time.Now().UTC(),
		Error:   message,
		Details: err.Error(),
		Code:    status,
	}
	if logger != nil {
		logger.Errorw("HTTP error", commonkeys.Error, err, "message", message, commonkeys.Status, status)
	}
	WriteJSON(w, status, body, headers...)
}

// WriteDecodeError is a shortcut for returning a decoding/binding error response.
func WriteDecodeError(w http.ResponseWriter, err error, logger output.ContextLogger, headers ...map[string]string) {
	WriteError(w, err, "Invalid request body", logger, headers...)
}

// WriteAuthError is a shortcut for returning an authentication/authorization error response.
func WriteAuthError(w http.ResponseWriter, err error, logger output.ContextLogger, headers ...map[string]string) {
	WriteError(w, err, "Unauthorized", logger, headers...)
}

// WriteNoContent responds with HTTP 204 No Content, optionally adding headers.
func WriteNoContent(w http.ResponseWriter, headers ...map[string]string) {
	for _, hdr := range headers {
		for k, v := range hdr {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
