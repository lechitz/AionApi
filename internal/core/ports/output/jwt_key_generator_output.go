package output

// JWTKeyGenerator defines an interface for generating JWT secret keys.
type JWTKeyGenerator interface {
	GenerateKey() (string, error)
}
