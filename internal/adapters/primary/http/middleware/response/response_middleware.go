// Package response provides common HTTP response handling functions and middleware.
package response

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Return sends an HTTP response with the specified status code and body, logging errors if writing the body fails.
func Return(w http.ResponseWriter, statusCode int, body []byte, logger output.Logger) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if len(body) != 0 {
		if _, err := w.Write(body); err != nil {
			logger.Errorw("failed to write response body", "error", err)
		}
	}
}

// ObjectResponse creates a JSON response with the given object, message, and current UTC date and returns it as a byte.Buffer.
func ObjectResponse(obj any, message string, logger output.Logger) *bytes.Buffer {
	response := struct {
		Date    time.Time `json:"date,omitempty"`
		Result  any       `json:"result,omitempty"`
		Message string    `json:"message,omitempty"`
	}{
		Message: message,
		Result:  obj,
		Date:    time.Now().UTC(),
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(response); err != nil {
		logger.Errorw("failed to encode response object to JSON", "error", err)
	}

	return body
}

// HandleError logs the error or warning, creates a JSON response, and sends it with the specified status code to the HTTP client.
func HandleError(w http.ResponseWriter, logger output.Logger, status int, msg string, err error) {
	if err != nil {
		logger.Errorw("operation failed",
			"message", msg,
			"error", err.Error(),
			"status", status,
		)
		response := ObjectResponse(nil, msg+": "+err.Error(), logger)
		Return(w, status, response.Bytes(), logger)
	} else {
		logger.Warnw("operation returned warning",
			"message", msg,
			"status", status,
		)
		response := ObjectResponse(nil, msg, logger)
		Return(w, status, response.Bytes(), logger)
	}
}

// HandleCriticalError logs a critical error and message, and then panics with the error or message provided.
func HandleCriticalError(logger output.Logger, message string, err error) {
	if err != nil {
		logger.Errorw("critical failure",
			"message", message,
			"error", err.Error(),
		)
		panic(err)
	}

	logger.Errorw("critical failure",
		"message", message,
	)
	panic(message)
}
