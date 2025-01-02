package input

import "github.com/lechitz/AionApi/internal/core/domain"

type ILoginService interface {
	GetUserByUsername(contextControl domain.ContextControl, user domain.LoginDomain) (domain.LoginDomain, error)
}
