package cache

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

type Creator interface {
	CreateToken(ctx context.Context, token domain.TokenDomain) (string, error)
}

type TokenChecker interface {
	Get(ctx context.Context, token domain.TokenDomain) (string, error)
}

type TokenSaver interface {
	Save(ctx context.Context, token domain.TokenDomain) error
}

type TokenUpdater interface {
	Update(ctx context.Context, token domain.TokenDomain) error
}

type TokenDeleter interface {
	Delete(ctx context.Context, token domain.TokenDomain) error
}

type TokenVerify interface {
	VerifyToken(ctx context.Context, token string) (uint64, string, error)
}

type TokenRepositoryPort interface {
	TokenChecker
	TokenSaver
	TokenUpdater
	TokenDeleter
}
