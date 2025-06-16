// Package constants contains constants related to authentication operations.
package constants

// ErrorUnauthorizedAccessMissingToken is returned when no authentication token is present in the request.
const ErrorUnauthorizedAccessMissingToken = "Unauthorized access: missing token"

// ErrorUnauthorizedAccessInvalidToken is returned when the authentication token provided is invalid.
const ErrorUnauthorizedAccessInvalidToken = "Unauthorized access: invalid token"

// Error is a generic error message.
const Error = "error"

// TokenKey is the key used for storing the token value.
const TokenKey = "token"

// UserIDKey is the key used to store the user ID.
const UserIDKey = "user_id"

// AuthToken is the key used to store the authentication token.
const AuthToken = "auth_token"
