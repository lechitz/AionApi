// Package constants contains constants related to authentication operations.
package constants

// ErrorUnauthorizedAccessMissingToken is returned when no authentication token is present in the request.
const ErrorUnauthorizedAccessMissingToken = "Unauthorized access: missing token"

// ErrorUnauthorizedAccessInvalidToken is returned when the authentication token provided is invalid.
const ErrorUnauthorizedAccessInvalidToken = "Unauthorized access: invalid token"
