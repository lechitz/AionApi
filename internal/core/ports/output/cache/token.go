package cache

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

type TokenCreator interface {
	CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
}

type TokenChecker interface {
	Check(ctx domain.ContextControl, token string) (uint64, string, error)
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

type TokenService interface {
	TokenCreator
	TokenChecker
	TokenSaver
	TokenUpdater
	TokenDeleter
}
