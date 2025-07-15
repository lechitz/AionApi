package output

// TokenProvider defines an interface for generating and parsing JWT tokens.
type TokenProvider interface {
	GenerateToken(userID uint64) (string, error)
	ParseToken(tokenString string) (map[string]interface{}, error)
}
