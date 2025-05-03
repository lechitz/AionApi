package constants

// User Constants

const (
	UserID = "user_id"
)

// Error User Messages

const (
	ErrorToDecodeUserRequest            = "error to decode user request"
	ErrorToCreateUser                   = "error to create user"
	ErrorToGetUser                      = "error to get user"
	ErrorToGetUsers                     = "error to get users"
	ErrorToUpdateUser                   = "error to update user"
	ErrorToSoftDeleteUser               = "error to soft delete user"
	ErrorToParseUser                    = "error to parse user"
	ErrorUnauthorizedAccessMissingToken = "error unauthorized access missing token"
)

// Success User Messages

const (
	SuccessToCreateUser     = "user created successfully"
	SuccessToGetUser        = "user get successfully"
	SuccessToGetUsers       = "users get successfully"
	SuccessToUpdateUser     = "user updated successfully"
	SuccessToUpdatePassword = "password updated successfully"
	SuccessUserSoftDeleted  = "user deleted successfully"
)
