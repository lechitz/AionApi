// Package ctxkeys contains constants used for context keys.
package ctxkeys

// ContextKey is a string type used for context keys.
type ContextKey string

// UserID is a constant string representing the key used to define or identify a user ID.
const UserID ContextKey = "user_id"

// Token is a constant string representing the key used to define or identify a token.
const Token ContextKey = "token"
