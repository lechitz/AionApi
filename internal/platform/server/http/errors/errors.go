// Package httperrors defines common sentinel errors used by HTTP handlers.
package httperrors

import "errors"

var (
	// ErrMethodNotAllowed represents a 405 response.
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrResourceNotFound represents a 404 response.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrInternalServer represents a generic internal server error.
	ErrInternalServer = errors.New("internal server error")
)
