// Package keygen defines an interface for generating JWT secret keys.
package keygen

// Generator defines an interface for generating JWT secret keys.
type Generator interface {
	Generate() (string, error)
}
