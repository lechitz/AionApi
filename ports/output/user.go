package output

import "github.com/lechitz/AionApi/internal/core/domain"

type IUserDomainDataBaseRepository interface {
	CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
}
