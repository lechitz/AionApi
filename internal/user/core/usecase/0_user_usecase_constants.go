// Package usecase constants contains constants related to user operations.
package usecase

import "errors"

// =============================================================================
// BUSINESS LOGIC - Roles
// =============================================================================

// UserRoles is the role of a user.
const UserRoles = "user"

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrorToHashPassword indicates an error while hashing a password.
	// #nosec G101: This constant does not leak a real secret, just an error message.
	ErrorToHashPassword = "error hashing password"

	// ErrorToCreateUser indicates an error when creating a user.
	ErrorToCreateUser = "error to create user"

	// ErrorToCompareHashAndPassword indicates a password hash comparison failure.
	ErrorToCompareHashAndPassword = "error to compare hash and password"

	// ErrorToCreateToken indicates an error when creating a token.
	ErrorToCreateToken = "error to create token"

	// ErrorToGetSelf indicates an error when fetching a user by ID.
	ErrorToGetSelf = "error to get user by id"

	// ErrorNoFieldsToUpdate indicates there were no fields to update in the user.
	ErrorNoFieldsToUpdate = "no fields to update"

	// ErrorToUpdatePassword indicates an error when updating the user password.
	ErrorToUpdatePassword = "error to update password"

	// ErrorToUpdateUser indicates an error when updating the user.
	ErrorToUpdateUser = "error to update user"

	// ErrorToGetUserByUsername indicates an error when fetching a user by username.
	ErrorToGetUserByUsername = "error to get user by username"

	// ErrorToSoftDeleteUser indicates an error when performing a soft delete on a user.
	ErrorToSoftDeleteUser = "error to soft delete user"
)

// =============================================================================
// BUSINESS LOGIC - Success Messages
// =============================================================================

const (
	// SuccessUserCreated indicates that the user was created successfully.
	SuccessUserCreated = "user created successfully"

	// SuccessUserRetrieved indicates a user was successfully retrieved.
	SuccessUserRetrieved = "user retrieved successfully"

	// SuccessPasswordUpdated indicates the password was updated successfully.
	SuccessPasswordUpdated = "password updated successfully"

	// SuccessUserUpdated indicates the user was updated successfully.
	SuccessUserUpdated = "user updated successfully"

	// SuccessUserSoftDeleted indicates a user was softly deleted successfully.
	SuccessUserSoftDeleted = "user soft deleted successfully"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrHashPassword is a sentinel error for password hashing failures.
	ErrHashPassword = errors.New(ErrorToHashPassword)

	// ErrCreateUser is a sentinel error for user creation failures.
	ErrCreateUser = errors.New(ErrorToCreateUser)

	// ErrCompareHashAndPassword is a sentinel error for password comparison failures.
	ErrCompareHashAndPassword = errors.New(ErrorToCompareHashAndPassword)

	// ErrCreateToken is a sentinel error for token creation failures.
	ErrCreateToken = errors.New(ErrorToCreateToken)

	// ErrGetSelf is a sentinel error for retrieving user by ID.
	ErrGetSelf = errors.New(ErrorToGetSelf)

	// ErrNoFieldsToUpdate is a sentinel error when no fields are provided for update.
	ErrNoFieldsToUpdate = errors.New(ErrorNoFieldsToUpdate)

	// ErrUpdatePassword is a sentinel error for password update failures.
	ErrUpdatePassword = errors.New(ErrorToUpdatePassword)

	// ErrUpdateUser is a sentinel error for user update failures.
	ErrUpdateUser = errors.New(ErrorToUpdateUser)

	// ErrGetUserByUsername is a sentinel error for retrieving user by username.
	ErrGetUserByUsername = errors.New(ErrorToGetUserByUsername)

	// ErrSoftDeleteUser is a sentinel error for soft delete failures.
	ErrSoftDeleteUser = errors.New(ErrorToSoftDeleteUser)
)
