package security

import (
	"github.com/lechitz/AionApi/core/domain"
)

type ITokenService interface {
	Create(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
	Check(ctx domain.ContextControl, token string) (uint64, string, error)
	Save(ctx domain.ContextControl, token domain.TokenDomain) error
	Update(ctx domain.ContextControl, token domain.TokenDomain) error
	Delete(ctx domain.ContextControl, token domain.TokenDomain) error
}
