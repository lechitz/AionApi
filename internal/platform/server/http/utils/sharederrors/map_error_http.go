// Package sharederrors provides error mapping utilities for handler-level HTTP response translation.
package sharederrors

import (
	"errors"
	"net/http"

	httperrors "github.com/lechitz/AionApi/internal/platform/server/http/errors"
)

// MapErrorToHTTPStatus maps domain and validation errors to the correct HTTP status code.
func MapErrorToHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// Typed errors (via errors.As)
	var ve *ValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest
	}

	var ue *UnauthorizedError
	if errors.As(err, &ue) {
		return http.StatusUnauthorized
	}

	var fe *ForbiddenError
	if errors.As(err, &fe) {
		return http.StatusForbidden
	}

	var mue *MissingUserIDError
	if errors.As(err, &mue) {
		return http.StatusUnauthorized
	}

	// Sentinel errors (via errors.Is) - include common HTTP handlers
	switch {
	case errors.Is(err, httperrors.ErrResourceNotFound):
		return http.StatusNotFound
	case errors.Is(err, httperrors.ErrMethodNotAllowed):
		return http.StatusMethodNotAllowed
	case errors.Is(err, ErrParseUserID):
		return http.StatusBadRequest
	case errors.Is(err, ErrNoFieldsToUpdate):
		return http.StatusBadRequest
	case errors.Is(err, ErrUsernameExists),
		errors.Is(err, ErrEmailExists),
		errors.Is(err, ErrDomainConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// ErrDomainConflict is used for business/domain conflicts (e.g., duplicate username/email).
var ErrDomainConflict = errors.New("domain conflict")
