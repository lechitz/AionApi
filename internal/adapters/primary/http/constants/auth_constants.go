package messages

// Auth Constants

const (
	AuthToken = "auth_token"
	Path      = "/"
	Domain    = "localhost"
)

// Error Auth Messages

const (
	ErrorToDecodeLoginRequest = "error to decode login request"
	ErrorToLogin              = "error to login"
	ErrorToRetrieveToken      = "error to retrieve token"
	ErrorToRetrieveUserID     = "error to retrieve user id"
	ErrorToLogout             = "error to logout"
)

// Success Auth Messages

const (
	SuccessToLogin  = "success to login"
	SuccessToLogout = "success to logout"
)
