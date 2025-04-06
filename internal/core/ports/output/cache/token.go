package cache

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

type Creator interface {
	CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
}

type TokenChecker interface {
	Get(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
}

type TokenSaver interface {
	Save(ctx domain.ContextControl, token domain.TokenDomain) error
}

type TokenUpdater interface {
	Update(ctx domain.ContextControl, token domain.TokenDomain) error
}

type TokenDeleter interface {
	Delete(ctx domain.ContextControl, token domain.TokenDomain) error
}

type TokenVerify interface {
	VerifyToken(ctx domain.ContextControl, token string) (uint64, string, error)
}

type TokenRepositoryPort interface {
	TokenChecker
	TokenSaver
	TokenUpdater
	TokenDeleter
}
