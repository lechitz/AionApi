package security

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordAdapter struct{}

func (BcryptPasswordAdapter) HashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hashed), err
}

func (BcryptPasswordAdapter) ComparePasswords(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
