// Package output defines the output ports for the auth package.
package output

// AuthProvider defines an interface for generating and verifying tokens.
type AuthProvider interface {
	Generate(userID uint64) (string, error)
	Verify(tokenValue string) (map[string]any, error)
	GenerateWithClaims(userID uint64, extra map[string]any) (string, error)
}
