// Package dto (user) contains data transfer objects for the HTTP layer.
package dto

import (
	"time"

	"github.com/lechitz/AionApi/internal/user/core/ports/input"
)

// UpdateUserRequest represents the payload for updating user information.
// Any field left nil will NOT update the corresponding user attribute.
// At least one field must be provided.
type UpdateUserRequest struct {
	// Name is the new display name for the user.
	// Example: "Alice Doe"
	Name *string `json:"name,omitempty" example:"Alice Doe"`

	// Username is the new unique handle for the user.
	// Example: "alice"
	Username *string `json:"username,omitempty" example:"alice"`

	// Email is the new email address for the user.
	// Example: "alice@example.com"
	Email *string `json:"email,omitempty" example:"alice@example.com"`

	// Locale is the optional locale (e.g., "en-US").
	Locale *string `json:"locale,omitempty" example:"en-US"`

	// Timezone is the optional timezone (IANA).
	Timezone *string `json:"timezone,omitempty" example:"America/Sao_Paulo"`

	// Location is an optional city/country description.
	Location *string `json:"location,omitempty" example:"São Paulo, BR"`

	// Bio is an optional short bio for the user.
	Bio *string `json:"bio,omitempty" example:"Backend engineer passionate about observability."`

	// AvatarURL is an optional URL for the user's avatar.
	AvatarURL *string `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`

	// OnboardingCompleted indicates if the user has completed the onboarding flow.
	OnboardingCompleted *bool `json:"onboarding_completed,omitempty" example:"true"`
}

// Validate ensures provided fields meet basic constraints.
func (r UpdateUserRequest) Validate() error {
	return validateOptionalProfileFields(r.Locale, r.Timezone, r.Location, r.Bio, r.AvatarURL)
}

// ToCommand converts the request to a domain command.
func (r UpdateUserRequest) ToCommand() input.UpdateUserCommand {
	return input.UpdateUserCommand{
		Name:                normalizeOptional(r.Name),
		Username:            normalizeOptional(r.Username),
		Email:               normalizeOptional(r.Email),
		Locale:              normalizeOptional(r.Locale),
		Timezone:            normalizeOptional(r.Timezone),
		Location:            normalizeOptional(r.Location),
		Bio:                 normalizeOptional(r.Bio),
		AvatarURL:           normalizeOptional(r.AvatarURL),
		OnboardingCompleted: r.OnboardingCompleted,
	}
}

// UpdateUserResponse represents the response returned after a successful user update.
type UpdateUserResponse struct {
	// UpdatedAt is the timestamp when the user was updated.
	// Example: "2025-09-14T22:01:02Z"
	UpdatedAt time.Time `json:"updated_at" format:"date-time" example:"2025-09-14T22:01:02Z"`

	// Name is the current display name after the update (if changed).
	// Example: "Alice Doe"
	Name *string `json:"name" example:"Alice Doe"`

	// Username is the current username after the update (if changed).
	// Example: "alice"
	Username *string `json:"username" example:"alice"`

	// Email is the current email after the update (if changed).
	// Example: "alice@example.com"
	Email *string `json:"email" example:"alice@example.com"`

	// ID is the user's unique identifier.
	// Example: 42
	ID uint64 `json:"user_id" example:"42"`

	// Locale returned after update.
	Locale *string `json:"locale,omitempty" example:"en-US"`

	// Timezone returned after update.
	Timezone *string `json:"timezone,omitempty" example:"America/Sao_Paulo"`

	// Location returned after update.
	Location *string `json:"location,omitempty" example:"São Paulo, BR"`

	// Bio returned after update.
	Bio *string `json:"bio,omitempty" example:"Backend engineer passionate about observability."`

	// AvatarURL returned after update.
	AvatarURL *string `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`

	// OnboardingCompleted indicates if the user has completed the onboarding flow.
	OnboardingCompleted bool `json:"onboarding_completed" example:"true"`
}
