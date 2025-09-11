// Package sharederrors centralizes application-wide error definitions and helpers.
package sharederrors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// ====================================================================
//                    ─── ERROR MESSAGE CONSTANTS.
// ====================================================================

const (
	// ErrMsgMissingUserID is the error message for missing user ID in context.
	ErrMsgMissingUserID = "missing user id in context"

	// ErrMsgMissingField is the error message for missing required fields.
	ErrMsgMissingField = "missing field: %s"

	// ErrMsgMissingFields is the error message for missing required fields.
	ErrMsgMissingFields = "required fields missing: %s"

	// ErrMsgUnauthorized is the error message for unauthorized access attempts.
	ErrMsgUnauthorized = "unauthorized"

	// ErrMsgForbidden is the error message for forbidden access attempts.
	ErrMsgForbidden = "forbidden"

	// ErrMsgValidation is the error message for validation failures.
	ErrMsgValidation = "validation error"

	// ErrMsgParseUserID is the error message for parsing user ID.
	ErrMsgParseUserID = "error parsing user id"

	// ErrUsernameInUse is the error message for a username already in use.
	ErrUsernameInUse = "username already in use"

	// ErrEmailInUse is the error message for an email already in use.
	ErrEmailInUse = "email already in use"

	// ErrUsernameAndEmailInUse is the error message for a username and email already in use.
	ErrUsernameAndEmailInUse = "username and email already in use"

	// ErrMsgCreateToken indicates an error when creating a token.
	ErrMsgCreateToken = "error to create token"

	// ErrMsgDeleteToken indicates an error when deleting a token.
	ErrMsgDeleteToken = "error to delete token"

	// NoFieldsToUpdate indicates an error when updating a user.
	NoFieldsToUpdate = "no fields to update"

	// ErrInvalidToken indicates an error when a token is invalid.
	ErrInvalidToken = "invalid token"

	// ErrTokenNotFound indicates an error when a token is not found.
	ErrTokenNotFound = "token not found"

	// ErrInvalidClaims indicates an error when a token is invalid.
	ErrInvalidClaims = "invalid claimsextractor"
)

// ErrNoFieldsToUpdate indicates an error when updating a user.
var ErrNoFieldsToUpdate = errors.New("no fields provided for update")

// ====================================================================
//                  ─── CUSTOM ERROR TYPES.
// ====================================================================

// MissingUserIDError represents an error when a user ID is missing from context.
type MissingUserIDError struct{}

func (e *MissingUserIDError) Error() string {
	return ErrMsgMissingUserID
}

// UnauthorizedError describes an error for unauthorized access attempts.
type UnauthorizedError struct {
	Reason string
}

func (e *UnauthorizedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s", ErrMsgUnauthorized, e.Reason)
	}
	return ErrMsgUnauthorized
}

// ForbiddenError describes an error for forbidden access attempts (authenticated, mas sem permissão).
type ForbiddenError struct {
	Reason string
}

func (e *ForbiddenError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s", ErrMsgForbidden, e.Reason)
	}
	return ErrMsgForbidden
}

// ValidationError represents a validation failure, including the affected field and reason.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	if e.Field != "" && e.Reason != "" {
		return fmt.Sprintf("%s on %s: %s", ErrMsgValidation, e.Field, e.Reason)
	}
	if e.Field != "" {
		return fmt.Sprintf("%s on %s", ErrMsgValidation, e.Field)
	}
	return ErrMsgValidation
}

// ====================================================================
//                ─── ERROR CONSTRUCTOR HELPERS.
// ====================================================================

// ErrMissingUserID returns a new MissingUserIDError.
func ErrMissingUserID() error {
	return &MissingUserIDError{}
}

// ErrUnauthorized returns a new UnauthorizedError with a specific reason.
func ErrUnauthorized(reason string) error {
	return &UnauthorizedError{Reason: reason}
}

// ErrForbidden returns a new ForbiddenError with a specific reason.
func ErrForbidden(reason string) error {
	return &ForbiddenError{Reason: reason}
}

// NewValidationError returns a new ValidationError for the given field and reason.
func NewValidationError(field, reason string) error {
	return &ValidationError{Field: field, Reason: reason}
}

// ====================================================================
//             ─── SENTINEL ERRORS (for errors.Is checks).
// ====================================================================

// ErrParseUserID is returned when a user ID cannot be parsed.
var ErrParseUserID = errors.New(ErrMsgParseUserID)

var ErrUsernameExists = errors.New(ErrUsernameInUse)
var ErrEmailExists = errors.New(ErrEmailInUse)

// ====================================================================
//             ─── ERRORS.
// ====================================================================

// AtLeastOneFieldRequired returns a new error with a list of required fields.
func AtLeastOneFieldRequired(fields ...string) error {
	return fmt.Errorf("at least one of the following fields must be provided: %s", strings.Join(fields, ", "))
}

// MissingFields returns a new error with a list of missing fields.
// If only one field is missing, the error message will be formatted as ErrMsgMissingField.
// If multiple fields are missing, the error message will be formatted as ErrMsgMissingFields.
// If no fields are missing, the error will be nil.
func MissingFields(fields ...string) error {
	if len(fields) == 1 {
		return NewValidationError(fields[0], fmt.Sprintf(ErrMsgMissingField, fields[0]))
	}
	return NewValidationError(commonkeys.Fields, fmt.Sprintf(ErrMsgMissingFields, strings.Join(fields, ", ")))
}
