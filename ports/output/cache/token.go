package cache

import (
	"github.com/lechitz/AionApi/core/domain/entities"
)

type ITokenRepository interface {
	SaveToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error
	GetToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) (string, error)
	UpdateToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error
	DeleteToken(ctx entities.ContextControl, tokenDomain entities.TokenDomain) error
}
