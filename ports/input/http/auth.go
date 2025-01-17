package http

import (
	"github.com/lechitz/AionApi/core/domain/entities"
)

type IAuthService interface {
	Login(ctx entities.ContextControl, userDomain entities.UserDomain, passwordReq string) (entities.UserDomain, string, error)
	Logout(ctx entities.ContextControl, token string) error
}
