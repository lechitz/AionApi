// Package constants contains constants related to authentication operations.
package constants

// Error is a generic error message.
const Error = "error"

// Username is the key used for usernames in authentication operations.
const Username = "username"

// ErrorToCompareHashAndPassword indicates an invalid password was provided.
const ErrorToCompareHashAndPassword = "invalid credentials" // #nosec G101

// ErrorToCreateToken indicates failure to create a token.
const ErrorToCreateToken = "error to create token"

// ErrorToCheckToken indicates failure to check a token.
const ErrorToCheckToken = "error to check token"

// ErrorToRevokeToken indicates failure to revoke a token.
const ErrorToRevokeToken = "error to revoke token"

// ErrorToGetUserByUserName indicates failure to retrieve a user by username.
const ErrorToGetUserByUserName = "error to get user by username"

// UserNotFoundOrInvalidCredentials indicates the user was not found or the provided credentials were invalid.
const UserNotFoundOrInvalidCredentials = "user not found or invalid credentials"

// InvalidCredentials indicates the provided credentials are invalid.
const InvalidCredentials = "invalid credentials" // #nosec G101

// SuccessToLogin indicates the user has logged in successfully.
const SuccessToLogin = "user logged in successfully"

// SuccessUserLoggedOut indicates the user has logged out successfully.
const SuccessUserLoggedOut = "user logged out successfully"

// Token is the key used for tokens in authentication operations.
const Token = "token"

// UserID is the key used to identify a user.
const UserID = "user_id"
