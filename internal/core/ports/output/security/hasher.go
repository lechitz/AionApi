// Package security provides interfaces for password hashing and validation.
package security

// Store defines methods for password hashing and validation.
// It provides mechanisms to securely hash and compare passwords.
type Store interface {
	HashPassword(plain string) (string, error)
	ValidatePassword(hashed, plain string) error
}
