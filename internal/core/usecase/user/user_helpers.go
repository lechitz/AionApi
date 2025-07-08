package user

import (
	"errors"
	"strings"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/badoux/checkmail"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// validateCreateUserRequired validates required fields for creating a user and returns an error if any validation fails.
func (s *Service) validateCreateUserRequired(user domain.UserDomain, password string) error {
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

// normalizeUserData adjusts user fields by trimming spaces, converting email to lowercase, and ensuring data uniformity. Returns the normalized user.
func (s *Service) normalizeUserData(user *domain.UserDomain) domain.UserDomain {
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
