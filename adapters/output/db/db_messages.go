package db

// Error messages for user

const (
	ErrorToCreateUser         = "error to create user: %w"
	ErrorToGetAllUsers        = "error to get all users: %w"
	ErrorToGetUserByID        = "error to get user by ID"
	ErrorToGetUserByUserName  = "error to get user by username"
	ErrorToUpdateUser         = "error to update user"
	ErrorToUpdateUserPassword = "error to update password"
	ErrorToSoftDeleteUser     = "error to soft delete user"
)

const (
	SuccesfullyDeletedUser = "Soft deleted user successfully"
)

const (
	DeleteAtIsNull = "deleted_at IS NULL"
)

// Table names

const (
	TableUsers = "aion_api.users"
)
