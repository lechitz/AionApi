// Package constants contains constants related to authentication operations.
package constants

// Error is a generic error message.
const Error = "error"

// ErrorToCompareHashAndPassword indicates an invalid password was provided.
const ErrorToCompareHashAndPassword = "invalid password provided"

// ErrorToCreateToken indicates failure to create a token.
const ErrorToCreateToken = "error to create token"

// ErrorToCheckToken indicates failure to check a token.
const ErrorToCheckToken = "error to check token"

// ErrorToRevokeToken indicates failure to revoke a token.
const ErrorToRevokeToken = "error to revoke token"

// ErrorToGetUserByUserName indicates failure to retrieve a user by username.
const ErrorToGetUserByUserName = "error to get user by username"

// SuccessToLogin indicates the user has logged in successfully.
const SuccessToLogin = "user logged in successfully"

// SuccessUserLoggedOut indicates the user has logged out successfully.
const SuccessUserLoggedOut = "user logged out successfully"

// Token is the key used for tokens in authentication operations.
const Token = "token"

// UserID is the key used to identify a user.
const UserID = "user_id"

// Name is the key used for the user's name.
const Name = "name"

// Username is the key used for the user's username.
const Username = "username"

// Email is the key used for the user's email.
const Email = "email"

// Password is the key used for the user's password.
const Password = "password"

// UpdatedAt is the key used for the updated at timestamp.
const UpdatedAt = "updated_at"
