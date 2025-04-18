package constants

// Error messages for user

const (
	ErrorToCreateUser        = "error to create user: %w"
	ErrorToGetAllUsers       = "error to get all users: %w"
	ErrorToGetUserByID       = "error to get user by ID"
	ErrorToGetUserByUsername = "error to get user by username"
	ErrorToUpdateUser        = "error to update user"
	ErrorToGetUserByEmail    = "error to get user by email"
	ErrorToSoftDeleteUser    = "error to soft delete user"
)

// Table names

const (
	TableUsers = "aion_api.users"
)

const (
	Error     = "error"
	CreatedAt = "created_at"
	UpdatedAt = "updated_at"
	DeletedAt = "deleted_at"
)
