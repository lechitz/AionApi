// Package domain contains core business entities used throughout the application.
package domain

// Auth represents a token associated with a user in the system.
// It includes the token string and the corresponding user ID.
type Auth struct {
	Token string // JWT or session token string
	Key   uint64 // ID of the user to whom the token belongs
}
