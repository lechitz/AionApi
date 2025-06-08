package constants

// Auth constants used in cookie and session control.
const (
	Token     = "token"
	Path      = "/"
	Domain    = "localhost"
	AuthToken = "auth_token"
)

// Error messages used in the authentication flow.
const (
	ErrorToDecodeLoginRequest = "failed to decode login request payload"
	ErrorToLogin              = "authentication process failed"
	ErrorToRetrieveToken      = "unable to retrieve access reference"
	ErrorToRetrieveUserID     = "failed to extract user ID from request context"
	ErrorToLogout             = "error occurred during logout"
)

// Success messages returned by authentication operations.
const (
	SuccessLogin  = "login successful"
	SuccessLogout = "logout successful"
)
