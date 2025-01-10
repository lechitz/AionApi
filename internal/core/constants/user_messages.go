package constants

// Validation Messages for User
const (
	NameIsRequired     = "name is required"
	UserIdIsRequired   = "user ID is required"
	UsernameIsRequired = "username is required"
	EmailIsRequired    = "email is required"
	PasswordIsRequired = "password is required"
	InvalidEmail       = "invalid email"
	ErrorParsingUserId = "error parsing user ID"
)

// Creation Errors
const (
	ErrorToCreateUser         = "error creating user"
	ErrorToFormatCreateUser   = "error formatting user during creation"
	ErrorToValidateCreateUser = "error validating user during creation"
)

// Retrieval Errors
const (
	ErrorToGetAllUsers       = "error retrieving all users"
	ErrorToGetUserByID       = "error retrieving user by ID"
	ErrorToGetUserByUserName = "error retrieving user by username"
)

// Update Errors
const (
	ErrorToUpdateUser       = "error updating user"
	ErrorToFormatUpdateUser = "error formatting user during update"
)

// Deletion Errors
const (
	ErrorToSoftDeleteUser = "error soft deleting user"
)

// Success Messages
const (
	SuccessUserCreated     = "user created successfully"
	SuccessUsersRetrieved  = "users retrieved successfully"
	SuccessUserRetrieved   = "user retrieved successfully"
	SuccessUserUpdated     = "user updated successfully"
	SuccessUserSoftDeleted = "user soft deleted successfully"
)
