// Package def contains constants used for context keys.
package def

// ContextKey is a string type used for context keys.
type ContextKey string

// CtxUserID is a constant string representing the key used to define or identify a user ID.
const CtxUserID ContextKey = "user_id"

// CtxUsers is a constant string representing the key used to define or identify a list of users.
const CtxUsers ContextKey = "users"

// CtxUsername is a constant string representing the key used to define or identify a username.
const CtxUsername ContextKey = "username"

// CtxToken is a constant string representing the key used to define or identify a token.
const CtxToken ContextKey = "token"

// CtxCategory is a constant string representing the key used to define or identify a category.
const CtxCategory ContextKey = "category"

// CtxCategoryID is a constant string representing the key used to define or identify a category ID.
const CtxCategoryID ContextKey = "category_id"

// CtxCategoryName is a constant string representing the key used to define or identify a category name.
const CtxCategoryName ContextKey = "category_name"

// AuthToken is the key used to store the authentication token.
const AuthToken = "auth_token"

// XRequestID represents the HTTP header name for propagating request ID across services.
const XRequestID = "X-Request-ID"
