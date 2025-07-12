// Package sharederrors centralizes application-wide error definitions and helpers.
package sharederrors

import (
	"errors"
	"fmt"
)

// ====== CUSTOM ERROR TYPES ======

// MissingUserIDError represents an error when a user ID is missing from context.
type MissingUserIDError struct{}

// Error returns a descriptive error message for MissingUserIDError.
func (e *MissingUserIDError) Error() string {
	return "missing user id in context"
}

// UnauthorizedError describes an error for unauthorized access attempts.
type UnauthorizedError struct {
	Reason string
}

// Error returns a descriptive error message for UnauthorizedError.
func (e *UnauthorizedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("unauthorized: %s", e.Reason)
	}
	return "unauthorized"
}

// ValidationError represents a validation failure, including the affected field and reason.
type ValidationError struct {
	Field  string
	Reason string
}

// Error returns a descriptive error message for ValidationError.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Reason)
}

//
// ====== ERROR CONSTRUCTOR HELPERS ======
//

// ErrMissingUserID returns a new MissingUserIDError.
func ErrMissingUserID() error {
	return &MissingUserIDError{}
}

// ErrUnauthorized returns a new UnauthorizedError with a specific reason.
func ErrUnauthorized(reason string) error {
	return &UnauthorizedError{Reason: reason}
}

// NewValidationError returns a new ValidationError for the given field and reason.
func NewValidationError(field, reason string) error {
	return &ValidationError{Field: field, Reason: reason}
}

//
// ====== SENTINEL ERRORS (for errors.Is checks) ======
//

// ErrParseUserID is returned when a user ID cannot be parsed.
var ErrParseUserID = errors.New("error parsing user id")
