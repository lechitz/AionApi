// Package constants contains constants used for token operations.
package constants

// Key is the key used for storing token values in Redis.
const Key = "key"

// Error is a generic error string used as a key or placeholder.
const Error = "error"

// ErrorToSaveTokenToRedis indicates a failure to save a token in Redis.
const ErrorToSaveTokenToRedis = "error to save token to Redis"

// ErrorToUpdateTokenInRedis indicates a failure to update a token in Redis.
const ErrorToUpdateTokenInRedis = "error to update token in Redis"

// ErrorToGetTokenFromRedis indicates a failure to retrieve a token from Redis.
const ErrorToGetTokenFromRedis = "error to get token from Redis"

// ErrorToDeleteTokenFromRedis indicates a failure to delete a token from Redis.
const ErrorToDeleteTokenFromRedis = "error to delete token from Redis"

// SuccessToUpdateTokenInRedis indicates a successful update of a token in Redis.
const SuccessToUpdateTokenInRedis = "success to update token in Redis"

// SuccessToGetTokenFromRedis indicates successful retrieval of a token from Redis.
const SuccessToGetTokenFromRedis = "success to get token from Redis"
