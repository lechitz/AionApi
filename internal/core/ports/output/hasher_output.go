// Package output security provides interfaces for password hashing and validation.
package output

// HasherStore defines methods for password hashing and validation.
type HasherStore interface {
	HashPassword(plain string) (string, error)
	ValidatePassword(hashed, plain string) error
}
