package security

type IPasswordService interface {
	HashPassword(plain string) (string, error)
	ComparePasswords(hashed, plain string) error
}
