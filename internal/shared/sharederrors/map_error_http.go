// Package sharederrors provides error mapping utilities for handler-level HTTP response translation.
package sharederrors

import (
	"errors"
	"net/http"
)

// MapErrorToHTTPStatus maps domain and validation errors to the correct HTTP status code.
func MapErrorToHTTPStatus(err error) int {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest
	}
	var ue *UnauthorizedError
	if errors.As(err, &ue) {
		return http.StatusUnauthorized
	}
	var mue *MissingUserIDError
	if errors.As(err, &mue) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrDomainConflict) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

// ErrDomainConflict is used for business/domain conflicts (e.g., duplicate username/email).
var ErrDomainConflict = errors.New("domain conflict")
