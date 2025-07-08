// Package ctxkeys contains constants used for context keys.
package ctxkeys

// ContextKey is a string type used for context keys.
type ContextKey string

// UserID is a constant string representing the key used to define or identify a user ID.
const UserID ContextKey = "user_id"

// Users is a constant string representing the key used to define or identify a list of users.
const Users ContextKey = "users"

// Username is a constant string representing the key used to define or identify a username.
const Username ContextKey = "username"

// Token is a constant string representing the key used to define or identify a token.
const Token ContextKey = "token"

// Category is a constant string representing the key used to define or identify a category.
const Category ContextKey = "category"

// CategoryID is a constant string representing the key used to define or identify a category ID.
const CategoryID ContextKey = "category_id"

// CategoryName is a constant string representing the key used to define or identify a category name.
const CategoryName ContextKey = "category_name"

// AuthToken is the key used to store the authentication token.
const AuthToken ContextKey = "auth_token"

// XRequestID represents the HTTP header name for propagating request ID across services.
