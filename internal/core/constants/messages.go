package constants

const (

	// Errors related to user

	ErrorToCreateUser         = "error to create user"
	ErrorToGetUsers           = "error to get users"
	ErrorToUpdateUser         = "error to update user"
	ErrorToDeleteUser         = "error to delete user"
	ErrorToFormatUpdateUser   = "error to format update user"
	ErrorToFormatCreateUser   = "error to format create user"
	ErrorToValidateCreateUser = "error to validate user"
	ErrorToGetUserByID        = "error to get user by ID"
	ErrorToGetUserByUserName  = "error to get user by username"

	// Errors related to validation

	NameIsRequired     = "name is required"
	UsernameIsRequired = "username is required"
	EmailIsRequired    = "email is required"
	PasswordIsRequired = "password is required"
	InvalidEmail       = "invalid email"
)
