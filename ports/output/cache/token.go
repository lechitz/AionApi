package cache

import (
	"github.com/lechitz/AionApi/core/domain"
)

type ITokenRepository interface {
	SaveToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error
	GetToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) (string, error)
	UpdateToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error
	DeleteToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) error
}
