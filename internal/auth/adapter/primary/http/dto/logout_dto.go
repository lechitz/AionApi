package dto

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
