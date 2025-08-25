package output

// KeyGenerator defines an interface for generating JWT secret keys.
type KeyGenerator interface {
	Generate() (string, error)
}
