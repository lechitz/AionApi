package security

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordAdapter struct{}

func NewBcryptPasswordAdapter() BcryptPasswordAdapter {
	return BcryptPasswordAdapter{}
}

func (BcryptPasswordAdapter) HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (BcryptPasswordAdapter) ValidatePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
