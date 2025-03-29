package cache

import (
	"github.com/lechitz/AionApi/core/domain"
)

type ITokenRepository interface {
	Save(ctx domain.ContextControl, token domain.TokenDomain) error
	Get(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
	Update(ctx domain.ContextControl, token domain.TokenDomain) error
	Delete(ctx domain.ContextControl, token domain.TokenDomain) error
}
