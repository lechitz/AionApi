package constants

// Errors related to user
const (
	ErrorToCreateUser         = "error to create user"
	ErrorToDecodeUserRequest  = "error to decode user request"
	ErrorToCreateToken        = "error to create token"
	ErrorToParseUser          = "error to parse user"
	ErrorToGetUsers           = "error to get users"
	ErrorToGetUser            = "error to get user"
	ErrorToGetUserByUsername  = "error to get user by username"
	ErrorToExtractUserID      = "error to extract user ID"
	ErrorToUpdateUser         = "error to update user"
	ErrorUserPermissionDenied = "user permission denied"
)

// Errors related to authentication
const (
	ErrorToLogin              = "error to login"
	ErrorToVerifyPassword     = "error to verify password"
	ErrorToDecodeLoginRequest = "error to decode login request"
)

// Errors related healthcheck and not found
const (
	ErrorToHealthCheck   = "error to health check"
	ErrorNotFoundHandler = "handler not found"
)

// Success messages
const (
	SuccessToCreateUser = "user created successfully"
	SuccessToGetUser    = "user get successfully"
	SuccessToGetUsers   = "users get successfully"
	SuccessToUpdateUser = "user updated successfully"
	SuccessToDeleteUser = "user deleted successfully"
	SuccessToLogin      = "success to login"
)
