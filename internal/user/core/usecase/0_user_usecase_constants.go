// Package usecase constants contains constants related to user operations.
package usecase

// UserRoles is the role of a user.
const UserRoles = "user"

// ErrorToHashPassword indicates an error while hashing a password.
// #nosec G101: This constant does not leak a real secret, just an error message.
const ErrorToHashPassword = "error hashing password"

// ErrorToCreateUser indicates an error when creating a user.
const ErrorToCreateUser = "error to create user"

// SuccessUserCreated indicates that the user was created successfully.
const SuccessUserCreated = "user created successfully"

// ErrorToCompareHashAndPassword indicates a password hash comparison failure.
const ErrorToCompareHashAndPassword = "error to compare hash and password"

// ErrorToCreateToken indicates an error when creating a token.
const ErrorToCreateToken = "error to create token"

// ErrorToGetSelf indicates an error when fetching a user by ID.
const ErrorToGetSelf = "error to get user by id"

// SuccessUserRetrieved indicates a user was successfully retrieved.
const SuccessUserRetrieved = "user retrieved successfully"

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
