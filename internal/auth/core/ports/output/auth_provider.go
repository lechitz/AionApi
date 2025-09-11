package output

import "context"

// AuthProvider defines an interface for generating and verifying tokens.
type AuthProvider interface {
	Generate(ctx context.Context, userID uint64) (string, error)
	Verify(ctx context.Context, tokenValue string) (map[string]any, error)
	GenerateWithClaims(ctx context.Context, userID uint64, extra map[string]any) (string, error)
}
