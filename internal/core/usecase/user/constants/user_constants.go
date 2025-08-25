// Package constants contains constants related to user operations.
package constants

// ErrorToValidateCreateUser indicates an error during user creation validation.
const ErrorToValidateCreateUser = "validation error on Create"

// ErrorToHashPassword indicates an error while hashing a password.
// #nosec G101: This constant does not leak a real secret, just an error message.
const ErrorToHashPassword = "error hashing password"

// ErrorToCreateUser indicates an error when creating a user.
const ErrorToCreateUser = "error to create user"

// SuccessUserCreated indicates that the user was created successfully.
const SuccessUserCreated = "user created successfully"

// DBErrorCheckingUsername indicates an error when checking the username in the database.
const DBErrorCheckingUsername = "error checking username in database"

// DBErrorCheckingEmail indicates an error when checking the email in the database.
const DBErrorCheckingEmail = "error checking email in database"

// ErrorToGetAllUsers indicates an error when fetching all users.
const ErrorToGetAllUsers = "error to get all users"

// ErrorToCompareHashAndPassword indicates a password hash comparison failure.
const ErrorToCompareHashAndPassword = "error to compare hash and password"

// ErrorToCreateToken indicates an error when creating a token.
const ErrorToCreateToken = "error to create token"

// ErrorToSaveToken indicates an error when saving a token.
const ErrorToSaveToken = "error to save token"

// SuccessUsersRetrieved indicates users were successfully retrieved.
const SuccessUsersRetrieved = "users retrieved successfully"

// ErrorToGetSelf indicates an error when fetching a user by ID.
const ErrorToGetSelf = "error to get user by id"

// SuccessUserRetrieved indicates a user was successfully retrieved.
const SuccessUserRetrieved = "user retrieved successfully"

// ErrorToGetUserByEmail indicates an error when fetching a user by email.
const ErrorToGetUserByEmail = "error to get user by email"

// ErrorNoFieldsToUpdate indicates there were no fields to update in the user.
const ErrorNoFieldsToUpdate = "no fields to update"

// ErrorToUpdatePassword indicates an error when updating the user password.
const ErrorToUpdatePassword = "error to update password"

// ErrorToUpdateUser indicates an error when updating the user.
const ErrorToUpdateUser = "error to update user"

// SuccessPasswordUpdated indicates the password was updated successfully.
const SuccessPasswordUpdated = "password updated successfully"

// SuccessUserUpdated indicates the user was updated successfully.
const SuccessUserUpdated = "user updated successfully"

// ErrorToGetUserByUsername indicates an error when fetching a user by username.
const ErrorToGetUserByUsername = "error to get user by username"

// ErrorToSoftDeleteUser indicates an error when performing a soft delete on a user.
const ErrorToSoftDeleteUser = "error to soft delete user"

// SuccessUserSoftDeleted indicates a user was softly deleted successfully.
const SuccessUserSoftDeleted = "user soft deleted successfully"

// TracerName is the name of the tracer used for the user use case.
const TracerName = "aionapi.user.create"

// TODO: Avaliar se as const abaixo podem ir para "commonkeys" que não seja ContextKeys !

// UpdatedAt is the key used for the last update timestamp.
const UpdatedAt = "updated_at"
