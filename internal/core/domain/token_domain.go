// Package domain contains core business entities used throughout the application.
package domain

// Token represents a token associated with a user in the system.
// It includes the token string and the corresponding user ID.
type Token struct {
	Value string // JWT or session token string
	Key   uint64 // ID of the user to whom the token belongs
}
