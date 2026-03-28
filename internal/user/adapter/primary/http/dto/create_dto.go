// internal/adapter/primary/http/controllers/user/dto/create.go

// Package dto provides Data Transfer Objects for the user HTTP layer.
package dto

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/lechitz/aion-api/internal/user/core/ports/input"
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
	// InvalidLocale is the error message for an invalid locale format.
	InvalidLocale = "locale must be 2-16 characters (e.g., en, en-US)"
	// InvalidTimezone is the error message for an invalid timezone format.
	InvalidTimezone = "timezone must be up to 64 characters (e.g., America/Sao_Paulo)"
	// InvalidLocation is the error message for an invalid location length.
	InvalidLocation = "location must be up to 255 characters"
	// InvalidBio is the error message for an invalid bio length.
	InvalidBio = "bio must be up to 1000 characters"
	// InvalidAvatarURL is the error message for an invalid avatar URL.
	InvalidAvatarURL = "avatar_url must be a valid URL"

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

	// Locale is the optional user locale (e.g., "en-US").
	// Example: "en-US"
	Locale *string `json:"locale,omitempty" example:"en-US"`

	// Timezone is the optional user timezone (IANA format).
	// Example: "America/Sao_Paulo"
	Timezone *string `json:"timezone,omitempty" example:"America/Sao_Paulo"`

	// Location is the optional user location.
	// Example: "São Paulo, BR"
	Location *string `json:"location,omitempty" example:"São Paulo, BR"`

	// Bio is an optional short bio for the user.
	// Example: "Backend engineer passionate about observability."
	Bio *string `json:"bio,omitempty" example:"Backend engineer passionate about observability."`

	// AvatarURL is an optional URL for the user's avatar.
	// Example: "https://example.com/avatar.png"
	AvatarURL *string `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`
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

	if err := validateOptionalProfileFields(r.Locale, r.Timezone, r.Location, r.Bio, r.AvatarURL); err != nil {
		return err
	}

	return nil
}

// ToCommand converts the request into a domain command for user creation.
func (r *CreateUserRequest) ToCommand() input.CreateUserCommand {
	return input.CreateUserCommand{
		Name:      strings.TrimSpace(r.Name),
		Username:  strings.TrimSpace(r.Username),
		Email:     strings.TrimSpace(r.Email),
		Password:  r.Password,
		Locale:    normalizeOptional(r.Locale),
		Timezone:  normalizeOptional(r.Timezone),
		Location:  normalizeOptional(r.Location),
		Bio:       normalizeOptional(r.Bio),
		AvatarURL: normalizeOptional(r.AvatarURL),
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

func validateOptionalProfileFields(locale, timezone, location, bio, avatarURL *string) error {
	if locale != nil {
		val := strings.TrimSpace(*locale)
		if len(val) < 2 || len(val) > 16 {
			return errors.New(InvalidLocale)
		}
	}
	if timezone != nil {
		val := strings.TrimSpace(*timezone)
		if val == "" || len(val) > 64 {
			return errors.New(InvalidTimezone)
		}
	}
	if location != nil {
		val := strings.TrimSpace(*location)
		if len(val) > 255 {
			return errors.New(InvalidLocation)
		}
	}
	if bio != nil {
		val := strings.TrimSpace(*bio)
		if len(val) > 1000 {
			return errors.New(InvalidBio)
		}
	}
	if avatarURL != nil {
		val := strings.TrimSpace(*avatarURL)
		if val != "" {
			if _, err := url.ParseRequestURI(val); err != nil {
				return errors.New(InvalidAvatarURL)
			}
		}
	}
	return nil
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
