package http

import (
	domain2 "github.com/lechitz/AionApi/internal/core/domain"
)

type IUserService interface {
	CreateUser(ctx domain2.ContextControl, user domain2.UserDomain, password string) (domain2.UserDomain, error)
	GetAllUsers(ctx domain2.ContextControl) ([]domain2.UserDomain, error)
	GetUserByID(ctx domain2.ContextControl, id uint64) (domain2.UserDomain, error)
	GetUserByUsername(ctx domain2.ContextControl, username string) (domain2.UserDomain, error)
	UpdateUser(ctx domain2.ContextControl, user domain2.UserDomain) (domain2.UserDomain, error)
	UpdateUserPassword(ctx domain2.ContextControl, user domain2.UserDomain, oldPassword, newPassword string) (domain2.UserDomain, string, error)
	SoftDeleteUser(ctx domain2.ContextControl, id uint64) error
}
