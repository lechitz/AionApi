// Package output defines the output ports for the auth package.
package output

// AuthProvider defines an interface for generating and verifying tokens.
type AuthProvider interface {
	GenerateRefreshToken(userID uint64) (string, error) // TODO: Para Refresh Token.
	Verify(tokenValue string) (map[string]any, error)
	GenerateAccessToken(userID uint64, extra map[string]any) (string, error) // TODO: Para Access Token.
}
