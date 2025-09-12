// Package commonkeys defines shared keys for user fields in config, logs, and context.
package commonkeys

const (

	// User is the key for user value in configs, logging, or context.
	User = "user"

	// UserID is the key for identifying a user ID.
	UserID = "user_id"

	// Username is the key for identifying a username.
	Username = "username"

	// Name is the key for the user's name field (DB or struct).
	Name = "name"

	// Email is the key for the user's email.
	Email = "email"

	// Roles is the key for the user's role.
	Roles = "roles"

	// Password is the key for the user's password.
	Password = "password"

	// NewPassword is the key for the user's new password.
	NewPassword = "new_password"

	// UserCreatedAt is the key for the user's created_at field (legacy/compat).
	UserCreatedAt = "created_at"

	// UserUpdatedAt is the key for the user's updated_at field (legacy/compat).
	UserUpdatedAt = "updated_at"

	// UserDeletedAt is the key for the user's deleted_at field (legacy/compat).
	UserDeletedAt = "deleted_at"

	// Users is the key for a collection of users.
	Users = "users"

	// UsersCount is the key for the total count of users.
	UsersCount = "users_count"

	// UserUpdatedFields is the key for fields updated in a user record.
	UserUpdatedFields = "updated_fields"

	// UpdatedUsername is the key for a user's updated username.
	UpdatedUsername = "updated_username"
)
