package security

type Hasher interface {
	HashPassword(plain string) (string, error)
	ComparePasswords(hashed, plain string) error
}
