// Package constants contains constants related to authentication operations.
package constants

// ErrorUnauthorizedAccessMissingToken is returned when no authentication token is present in the request.
const ErrorUnauthorizedAccessMissingToken = "unauthorized access: missing token"

// ErrorUnauthorizedAccessInvalidToken is returned when the authentication token provided is invalid.
const ErrorUnauthorizedAccessInvalidToken = "unauthorized access: invalid token"
