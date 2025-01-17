package messages

// Error Auth Messages

const (
	ErrorToDecodeLoginRequest = "error to decode login request"
	ErrorToLogin              = "error to login"
	ErrorToRetrieveToken      = "error to retrieve token"
	ErrorToRetrieveUserID     = "error to retrieve user id"
	ErrorToLogout             = "error to logout"
	ErrorUserIDIsNil          = "userID is nil"
)

// Success Auth Messages

const (
	SuccessToLogin  = "success to login"
	SuccessToLogout = "success to logout"
)
