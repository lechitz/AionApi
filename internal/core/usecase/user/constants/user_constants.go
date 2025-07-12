// Package constants contains constants related to user operations.
package constants

// ErrorToValidateCreateUser indicates an error during user creation validation.
const ErrorToValidateCreateUser = "validation error on CreateUser"

// ErrorToDeleteToken indicates an error when deleting a token.
const ErrorToDeleteToken = "error to delete token"

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

// ErrorToGetUserByID indicates an error when fetching a user by ID.
const ErrorToGetUserByID = "error to get user by id"

// SuccessUserRetrieved indicates a user was successfully retrieved.
const SuccessUserRetrieved = "user retrieved successfully"

// ErrorToGetUserByUserName indicates an error when fetching a user by username.
const ErrorToGetUserByUserName = "error to get user by username"

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

// ErrorToSoftDeleteUser indicates an error when performing a soft delete on a user.
const ErrorToSoftDeleteUser = "error to soft delete user"

// SuccessUserSoftDeleted indicates a user was softly deleted successfully.
const SuccessUserSoftDeleted = "user soft deleted successfully"

// NameIsRequired indicates that the user's name is required.
const NameIsRequired = "name is required"

// UsernameIsRequired indicates that the user's username is required.
const UsernameIsRequired = "username is required"

// EmailIsRequired indicates that the user's email is required.
const EmailIsRequired = "email is required"

// PasswordIsRequired indicates that the user's password is required.
const PasswordIsRequired = "password is required"

// InvalidEmail indicates that the email format is invalid.
const InvalidEmail = "invalid email format"

// UsernameIsAlreadyInUse indicates the username is already taken.
const UsernameIsAlreadyInUse = "username is already in use"

// EmailIsAlreadyInUse indicates the email is already registered.
const EmailIsAlreadyInUse = "email is already in use"

// TracerName is the name of the tracer used for the user use case.
const TracerName = "aionapi.user.usecase"

// TODO: Avaliar se as const abaixo podem ir para "commonkeys" que não seja ContextKeys !

// UpdatedAt is the key used for the last update timestamp.
const UpdatedAt = "updated_at"
