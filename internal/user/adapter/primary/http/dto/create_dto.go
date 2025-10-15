// internal/adapter/primary/http/controllers/user/dto/create.go

// Package dto provides Data Transfer Objects for the user HTTP layer.
package dto

import (
	"errors"
	"regexp"

	"github.com/lechitz/AionApi/internal/user/core/ports/input"
)

const (
	// NameIsRequired is the error message for a missing name.
	NameIsRequired = "name is required"
	// UsernameIsRequired is the error message for a missing username.
	UsernameIsRequired = "username is required"
	// EmailIsRequired is the error message for a missing email.
	EmailIsRequired = "email is required"
	// PasswordIsRequired is the error message for a missing password.
	PasswordIsRequired = "password is required"
	// InvalidPasswordLength is the error message for an invalid password length.
	InvalidPasswordLength = "password must be at least 8 characters"
	// InvalidEmail is the error message for an invalid email format.
	InvalidEmail = "invalid email format"

	// emailRegex is the regular expression for email validation.
	emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

// CreateUserRequest represents the request payload for creating a user.
// Fields include examples to improve Swagger UI readability.
type CreateUserRequest struct {
	// Name is the human-friendly display name of the user.
	// Example: "João Pereira"
	Name string `json:"name" example:"João Pereira"`

	// Username is the unique handle for login and identification.
	// Example: "lechitz"
	Username string `json:"username" example:"lechitz"`

	// Email is the user's contact email address.
	// Example: "dev@aion.local"
	Email string `json:"email" example:"dev@aion.local"`

	// Password is the user's credential (minimum length: 8).
	// Example: "P@ssw0rd123"
	Password string `json:"password" example:"P@ssw0rd123"`
}

// ValidateUser validates the user input for required fields and basic constraints.
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

// ToCommand converts the request into a domain command for user creation.
func (r *CreateUserRequest) ToCommand() input.CreateUserCommand {
	return input.CreateUserCommand{
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

// CreateUserResponse is the response payload returned after a successful user creation.
type CreateUserResponse struct {
	// Name is the created user's display name.
	// Example: "João Pereira"
	Name string `json:"name" example:"João Pereira"`

	// Username is the created user's unique handle.
	// Example: "lechitz"
	Username string `json:"username" example:"lechitz"`

	// Email is the created user's email address.
	// Example: "dev@aion.local"
	Email string `json:"email" example:"dev@aion.local"`

	// ID is the created user's identifier.
	// Example: 42
	ID uint64 `json:"id" example:"42"`
}
