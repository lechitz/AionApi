// internal/adapters/primary/http/controllers/user/dto/create.go

package dto

import (
	"errors"
	"regexp"

	"github.com/lechitz/AionApi/internal/core/ports/input"
)

const (
	NameIsRequired        = "name is required"
	UsernameIsRequired    = "username is required"
	EmailIsRequired       = "email is required"
	PasswordIsRequired    = "password is required"
	InvalidPasswordLength = "password must be at least 8 characters"
	InvalidEmail          = "invalid email format"
	emailRegex            = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

// CreateUserRequest represents the request for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ValidateUser validates the user input.
func (r *CreateUserRequest) ValidateUser() error {
	if r.Name == "" {
		return errors.New(NameIsRequired)
	}
	if r.Username == "" {
		return errors.New(UsernameIsRequired)
	}
	if r.Email == "" {
		return errors.New(EmailIsRequired)
	}
	if r.Password == "" {
		return errors.New(PasswordIsRequired)
	}
	if len(r.Password) < 8 {
		return errors.New(InvalidPasswordLength)
	}
	if !regexp.MustCompile(emailRegex).MatchString(r.Email) {
		return errors.New(InvalidEmail)
	}
	return nil
}

// ToCommand converts the request to a command.
func (r *CreateUserRequest) ToCommand() input.CreateUserCommand {
	return input.CreateUserCommand{
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}
