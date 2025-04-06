package security

type SecurityStore interface {
	HashPassword(plain string) (string, error)
	ValidatePassword(hashed, plain string) error
}
