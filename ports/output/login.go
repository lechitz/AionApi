package output

import "github.com/lechitz/AionApi/internal/core/domain"

type ILoginDomainDataBaseRepository interface {
	GetUserByUsername(contextControl domain.ContextControl, user domain.LoginDomain) (domain.LoginDomain, error)
}
