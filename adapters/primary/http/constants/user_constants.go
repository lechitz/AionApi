package constants

// UserID is the key used to identify a user in context or request scope.
const UserID = "user_id"

// Error messages returned by user operations.
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

// Success messages returned by user operations.
const (
	SuccessToCreateUser     = "user created successfully"
	SuccessToGetUser        = "user get successfully"
	SuccessToGetUsers       = "users get successfully"
	SuccessToUpdateUser     = "user updated successfully"
	SuccessToUpdatePassword = "password updated successfully"
	SuccessUserSoftDeleted  = "user deleted successfully"
)
