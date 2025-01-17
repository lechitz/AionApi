package msg

// Validation Messages for User
const (
	NameIsRequired     = "name is required"
	UserIdIsRequired   = "user ID is required"
	UsernameIsRequired = "username is required"
	EmailIsRequired    = "email is required"
	PasswordIsRequired = "password is required"
	InvalidEmail       = "invalid email"
)

// Error Messages for User
const (
	ErrorToCreateUser             = "error creating user"
	ErrorToNormalizeUserData      = "error formatting user during creation"
	ErrorToValidateCreateUser     = "error validating user during creation"
	ErrorToGetAllUsers            = "error retrieving all users"
	ErrorToGetUserByID            = "error retrieving user by ID"
	ErrorToGetUserByUserName      = "error retrieving user by username"
	ErrorParsingUserId            = "error parsing user ID"
	ErrorToSoftDeleteUser         = "error soft deleting user"
	ErrorToUpdateUser             = "error updating user"
	ErrorToFormatUpdateUser       = "error formatting user during update"
	ErrorToHashPassword           = "error hashing password"
	ErrorToCompareHashAndPassword = "invalid password provided"
)

// Success Messages for User
const (
	SuccessUserCreated     = "user created successfully"
	SuccessUsersRetrieved  = "users retrieved successfully"
	SuccessUserRetrieved   = "user retrieved successfully"
	SuccessUserUpdated     = "user updated successfully"
	SuccessUserSoftDeleted = "user soft deleted successfully"
)
