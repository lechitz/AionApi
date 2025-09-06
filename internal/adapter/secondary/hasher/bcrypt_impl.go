// Package hasher provides a bcrypt-based implementation of output.Hasher.
package hasher

import "golang.org/x/crypto/bcrypt"

// BcryptHasher implements output.Hasher using bcrypt.
// It is stateless; use value receivers.
type BcryptHasher struct{}

// New returns a bcrypt hasher as a value (no external deps needed).
func New() BcryptHasher {
	return BcryptHasher{}
}

// Hash returns the bcrypt hash of the given plain text.
func (BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare verifies a bcrypt hash against plain text.
func (BcryptHasher) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
