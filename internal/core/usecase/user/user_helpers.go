package user

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"strings"
)

func (s *UserService) validateCreateUserRequired(user domain.UserDomain, password string) error {
	if user.Name == "" {
		return errors.New(constants.NameIsRequired)
	}
	if user.Username == "" {
		return errors.New(constants.UsernameIsRequired)
	}
	if user.Email == "" {
		return errors.New(constants.EmailIsRequired)
	}
	if password == "" {
		return errors.New(constants.PasswordIsRequired)
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New(constants.InvalidEmail)
	}
	return nil
}

func (s *UserService) normalizeUserData(user *domain.UserDomain) domain.UserDomain {
	if user.Name != "" {
		user.Name = strings.TrimSpace(user.Name)
	}
	if user.Username != "" {
		user.Username = strings.TrimSpace(user.Username)
	}
	if user.Email != "" {
		user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	}
	return *user
}
