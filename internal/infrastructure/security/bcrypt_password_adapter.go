package security

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordAdapter struct{}

func (BcryptPasswordAdapter) HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func (BcryptPasswordAdapter) ValidatePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
