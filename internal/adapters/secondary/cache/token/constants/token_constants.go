// Package constants contains constants used for token operations.
package constants

// Key is the key used for storing token values in Redis.
const Key = "key"

// Error is a generic error string used as a key or placeholder.
const Error = "error"

// UserID is the key used to store the user ID.
const UserID = "user_id"

// ErrorToSaveTokenToRedis indicates a failure to save a token in Redis.
const ErrorToSaveTokenToRedis = "error to save token to Redis"

// ErrorToUpdateTokenInRedis indicates a failure to update a token in Redis.
const ErrorToUpdateTokenInRedis = "error to update token in Redis"

// ErrorToGetTokenFromRedis indicates a failure to retrieve a token from Redis.
const ErrorToGetTokenFromRedis = "error to get token from Redis"

// ErrorToDeleteTokenFromRedis indicates a failure to delete a token from Redis.
const ErrorToDeleteTokenFromRedis = "error to delete token from Redis"
