package constants

const (
	ErrorToValidateCreateUser     = "error to validate create user"
	ErrorToDeleteToken            = "error to delete token"
	ErrorToHashPassword           = "error to hash password" // #nosec G101
	ErrorToCreateUser             = "error to create user"
	SuccessUserCreated            = "user created successfully"
	ErrorToGetAllUsers            = "error to get all users"
	ErrorToCompareHashAndPassword = "error to compare hash and password"
	ErrorToCreateToken            = "error to create token"
	ErrorToSaveToken              = "error to save token"
	SuccessUsersRetrieved         = "users retrieved successfully"
	ErrorToGetUserByID            = "error to get user by id"
	SuccessUserRetrieved          = "user retrieved successfully"
	ErrorToGetUserByUserName      = "error to get user by username"
	ErrorToGetUserByEmail         = "error to get user by email"
	ErrorNoFieldsToUpdate         = "no fields to update"
	ErrorToUpdatePassword         = "error to update password"
	ErrorToUpdateUser             = "error to update user"
	SuccessPasswordUpdated        = "password updated successfully"
	SuccessUserUpdated            = "user updated successfully"
	ErrorToSoftDeleteUser         = "error to soft delete user"
	SuccessUserSoftDeleted        = "user soft deleted successfully"
	NameIsRequired                = "name is required"
	UsernameIsRequired            = "username is required"
	EmailIsRequired               = "email is required"
	PasswordIsRequired            = "password is required"
	InvalidEmail                  = "invalid email format"
	UsernameIsAlreadyInUse        = "username is already in use"
	EmailIsAlreadyInUse           = "email is already in use"
)

const (
	UserID    = "user_id"
	Users     = "users"
	Error     = "error"
	Name      = "name"
	Username  = "username"
	Email     = "email"
	Password  = "password"
	UpdatedAt = "updated_at"
)
