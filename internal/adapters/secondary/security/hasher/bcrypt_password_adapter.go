package hasher

import "golang.org/x/crypto/bcrypt"

// BcryptHasher provides methods for hashing and validating passwords using the bcrypt algorithm.
type BcryptHasher struct{}

// NewBcryptHasher creates a new instance of BcryptHasher for password hashing and validation.
func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{}
}

// Hash generates a bcrypt-hashed string from the provided plaintext password. Returns the hashed password or an error if hashing fails.
func (BcryptHasher) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare compares a bcrypt-hashed password with a plaintext password and returns an error if they do not match.
func (BcryptHasher) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
