package token

import "github.com/lechitz/AionApi/internal/core/domain"

type TokenUsecase interface {
	CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
	VerifyToken(ctx domain.ContextControl, token string) (uint64, string, error)
	Delete(ctx domain.ContextControl, token domain.TokenDomain) error
	Update(ctx domain.ContextControl, token domain.TokenDomain) error
	Save(ctx domain.ContextControl, token domain.TokenDomain) error
}
