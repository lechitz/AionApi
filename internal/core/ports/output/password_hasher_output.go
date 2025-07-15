// Package output security provides interfaces for password hashing and validation.
package output

// Hasher defines methods for password hashing and validation.
type Hasher interface {
	Hash(plain string) (string, error)
	Compare(hashed, plain string) error
}
