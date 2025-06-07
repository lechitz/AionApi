package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

type TokenUsecase interface {
	CreateToken(ctx context.Context, token domain.TokenDomain) (string, error)
	VerifyToken(ctx context.Context, token string) (uint64, string, error)
	Delete(ctx context.Context, token domain.TokenDomain) error
	Update(ctx context.Context, token domain.TokenDomain) error
	Save(ctx context.Context, token domain.TokenDomain) error
}
