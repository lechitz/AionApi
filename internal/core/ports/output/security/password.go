package security

type PasswordManager interface { //TODO criar o mock
	HashPassword(plain string) (string, error)
	ComparePasswords(hashed, plain string) error
}
