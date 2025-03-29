package db

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

type IUserRepository interface {
	CreateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error)
	GetUserByID(ctx domain.ContextControl, ID uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error)
	GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error)
	UpdateUser(ctx domain.ContextControl, userID uint64, fields map[string]interface{}) (domain.UserDomain, error)
	SoftDeleteUser(ctx domain.ContextControl, userID uint64) error
}
