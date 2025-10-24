// Package dto contains Data Transfer Objects used by the HTTP layer.
package dto

// LoginUserRequest is the request body for login.
// Fields include examples to improve Swagger UI readability.
type LoginUserRequest struct {
	// Username is the unique identifier used to authenticate the user.
	// Example: "joaopereira"
	Username string `json:"username" example:"joaopereira"`

	// Password is the credential paired with the username.
	// Example: "P@ssw0rd123"
	Password string `json:"password" example:"P@ssw0rd123"`
}

// LoginUserResponse is the response body returned on a successful login.
type LoginUserResponse struct {
	// Token is the JSON Web Token used for authenticating subsequent requests.
	// Example: "eyJhbGciOi..."
	Token string `json:"token" example:"eyJhbGciOi..."`

	// ID is the unique identifier for the authenticated user.
	// Example: 42
	ID uint64 `json:"id" example:"42"`

	// Name is a friendly display name for the authenticated user.
	// Example: "João Pereira"
	Name string `json:"name" example:"João Pereira"`

	// Roles are the permissions or roles assigned to the user.
	// Example: ["admin", "user"]
	Roles []string `json:"roles" example:"admin,user"`
}

// LogoutUserRequest is NOT used by the current logout endpoint.
// The endpoint is stateless and relies on the auth context/cookie.
// This struct remains as a placeholder for future explicit-body flows (if needed).
type LogoutUserRequest struct {
	// Token preview or opaque token for explicit invalidation (not used today).
	// Example: "eyJhbGciOi..." (truncated)
	Token string `json:"token,omitempty" example:"eyJhbGciOi..."`

	// UserID for explicit ownership checks (not used today).
	// Example: 42
	UserID uint64 `json:"user_id,omitempty" example:"42"`
}
