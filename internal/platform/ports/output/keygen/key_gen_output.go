// Package crypto defines an interface for generating JWT secret keys.
package keygen

// KeyGenerator defines an interface for generating JWT secret keys.
type KeyGenerator interface {
	Generate() (string, error)
}
