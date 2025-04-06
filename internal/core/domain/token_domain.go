package domain

type TokenConfig struct {
	SecretKey string
}

type TokenDomain struct {
	UserID uint64
	Token  string
}
