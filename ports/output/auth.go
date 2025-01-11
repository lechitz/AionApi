package output

import (
	"context"
)

type IAuthRepository interface {
	CreateToken(ctx context.Context, userID uint64) (string, error)
	ValidateToken(ctx context.Context, tokenFromCookie string) (string, uint64, error)
	GetTokenByUserID(ctx context.Context, userID uint64) (string, error)
	UpdateTokenByUserID(ctx context.Context, userID uint64) (string, error)
	DeleteTokenByUserID(ctx context.Context, userID uint64) error
}
