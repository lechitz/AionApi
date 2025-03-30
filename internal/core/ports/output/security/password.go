package security

type PasswordManager interface {
	HashPassword(plain string) (string, error)
	ComparePasswords(hashed, plain string) error
}
