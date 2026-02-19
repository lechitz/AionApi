// Package dto contains Data Transfer Objects used by the HTTP layer.
package dto

// RefreshResponse represents the response body after successful token refresh.
type RefreshResponse struct {
	// Token is the new JSON Web Token to be used for authenticating subsequent requests.
	// Example: "eyJhbGciOi..."
	Token string `json:"token" example:"eyJhbGciOi..."`
}
