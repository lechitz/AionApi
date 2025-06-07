// Package dto contains Data Transfer Objects used by the HTTP layer.
package dto

// LoginUserRequest represents the expected payload to log a user in.
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUserResponse represents **the** response payload after a successful login.
type LoginUserResponse struct {
	Username string `json:"username"`
}

// LogoutUserRequest represents the payload used to request a logout.
type LogoutUserRequest struct {
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}
