// Package constants contains constants used throughout the application.
package constants

// Token is the key for authentication tokens used in cookies and sessions.
const Token = "token"

// Path defines the default path used for cookies.
const Path = "/"

// Domain is the default domain used for authentication cookies.
const Domain = "localhost"

// AuthToken is the key for the primary authentication token.
const AuthToken = "auth_token"

// ErrorToDecodeLoginRequest is returned when decoding the login request fails.
const ErrorToDecodeLoginRequest = "failed to decode login request payload"

// ErrorToLogin is returned when the authentication process fails.
const ErrorToLogin = "authentication process failed"

// ErrorToRetrieveToken is returned when the access reference cannot be retrieved.
const ErrorToRetrieveToken = "unable to retrieve access reference" // #nosec G101

// ErrorToRetrieveUserID is returned when extracting the user ID from the context fails.
const ErrorToRetrieveUserID = "failed to extract user ID from request context"

// ErrorToLogout is returned when an error occurs during logout.
const ErrorToLogout = "error occurred during logout"

// SuccessLogin indicates a successful login operation.
const SuccessLogin = "login successful"

// SuccessLogout indicates a successful logout operation.
const SuccessLogout = "logout successful"
