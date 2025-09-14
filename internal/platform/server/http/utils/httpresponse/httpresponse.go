// Package httpresponse centralizes HTTP response helpers for all controllers/domains.
package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Internal constants used to avoid magic strings in this package.
const (
	// headerContentType is the HTTP header name for content type.
	headerContentType = "Content-Type"
	// contentTypeJSON is the JSON content type used for API responses.
	contentTypeJSON = "application/json"

	// logMsgHTTPError is the standardized log message for HTTP errors.
	logMsgHTTPError = "HTTP error"
	// logMsgDecodeError is the standardized log message for decode/binding errors.
	logMsgDecodeError = "decode error"
	// logMsgDomainError is the standardized log message for domain errors.
	logMsgDomainError = "domain error"

	// msgUnauthorized is the default error message for authentication/authorization failures.
	msgUnauthorized = "Unauthorized"
	// msgInvalidRequestBody is the default error message for malformed request bodies.
	msgInvalidRequestBody = "Invalid request body"
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
	w.Header().Set(headerContentType, contentTypeJSON)
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
func WriteError(w http.ResponseWriter, err error, message string, log logger.ContextLogger, headers ...map[string]string) {
	status := sharederrors.MapErrorToHTTPStatus(err)
	body := ResponseBody{
		Date:    time.Now().UTC(),
		Error:   message,
		Details: err.Error(),
		Code:    status,
	}
	if log != nil {
		log.Errorw(logMsgHTTPError, commonkeys.Error, err, commonkeys.Message, message, commonkeys.Status, status)
	}
	WriteJSON(w, status, body, headers...)
}

// WriteDecodeError is a shortcut for returning a decoding/binding error response.
func WriteDecodeError(w http.ResponseWriter, err error, log logger.ContextLogger, headers ...map[string]string) {
	WriteError(w, err, msgInvalidRequestBody, log, headers...)
}

// WriteAuthError is a shortcut for returning an authentication/authorization error response.
func WriteAuthError(w http.ResponseWriter, err error, log logger.ContextLogger, headers ...map[string]string) {
	WriteError(w, err, msgUnauthorized, log, headers...)
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

// ---------------------------
// Span-aware helper variants
// ---------------------------

// WriteAuthErrorSpan records trace/log metadata for an auth error, then writes the HTTP response.
func WriteAuthErrorSpan(ctx context.Context, w http.ResponseWriter, span trace.Span, log logger.ContextLogger) {
	err := sharederrors.ErrMissingUserID()
	span.RecordError(err)
	span.SetStatus(codes.Error, sharederrors.ErrMsgMissingUserID)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusUnauthorized))
	if log != nil {
		log.ErrorwCtx(ctx, sharederrors.ErrMsgMissingUserID, commonkeys.Error, err.Error())
	}
	WriteAuthError(w, err, log)
}

// WriteDecodeErrorSpan records trace/log metadata for a decode/binding error, then writes the HTTP response.
func WriteDecodeErrorSpan(ctx context.Context, w http.ResponseWriter, span trace.Span, err error, log logger.ContextLogger) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusBadRequest))
	if log != nil {
		log.ErrorwCtx(ctx, logMsgDecodeError, commonkeys.Error, err.Error())
	}
	WriteDecodeError(w, err, log)
}

// WriteDomainErrorSpan records trace/log metadata for a domain error, then writes the HTTP response.
func WriteDomainErrorSpan(ctx context.Context, w http.ResponseWriter, span trace.Span, err error, message string, log logger.ContextLogger) {
	statusCode := sharederrors.MapErrorToHTTPStatus(err)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, statusCode))
	if log != nil {
		log.ErrorwCtx(ctx, logMsgDomainError, commonkeys.Error, err.Error())
	}
	WriteError(w, err, message, log)
}
