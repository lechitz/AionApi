package security

import "golang.org/x/crypto/bcrypt"

// BcryptPasswordAdapter provides methods for hashing and validating passwords using the bcrypt algorithm.
type BcryptPasswordAdapter struct{}

// NewBcryptPasswordAdapter creates a new instance of BcryptPasswordAdapter for password hashing and validation.
func NewBcryptPasswordAdapter() BcryptPasswordAdapter {
	return BcryptPasswordAdapter{}
}

// HashPassword generates a bcrypt-hashed string from the provided plaintext password. Returns the hashed password or an error if hashing fails.
func (BcryptPasswordAdapter) HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ValidatePassword compares a bcrypt-hashed password with a plaintext password and returns an error if they do not match.
func (BcryptPasswordAdapter) ValidatePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
