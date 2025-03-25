package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	HashPassword(password string) (string, error)
}

type PasswordChecker interface {
	ComparePassword(hashed, password string) error
}

type PasswordUtil struct{}

func NewPasswordUtil() *PasswordUtil {
	return &PasswordUtil{}
}

func (p *PasswordUtil) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (p *PasswordUtil) ComparePassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
