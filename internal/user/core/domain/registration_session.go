package domain

import "time"

// RegistrationSessionStatus represents the lifecycle state of a registration session.
type RegistrationSessionStatus string

const (
	// RegistrationStatusPending indicates a session waiting for completion.
	RegistrationStatusPending RegistrationSessionStatus = "pending"
	// RegistrationStatusCompleted indicates a session successfully completed.
	RegistrationStatusCompleted RegistrationSessionStatus = "completed"
	// RegistrationStatusExpired indicates a session expired before completion.
	RegistrationStatusExpired RegistrationSessionStatus = "expired"
	// RegistrationStatusCanceled indicates a session canceled by user/system action.
	RegistrationStatusCanceled RegistrationSessionStatus = "canceled"
)

// RegistrationSession stores the in-progress multi-step public registration flow.
type RegistrationSession struct {
	RegistrationID string
	Name           string
	Username       string
	Email          string
	PasswordHash   string
	Locale         *string
	Timezone       *string
	Location       *string
	Bio            *string
	AvatarURL      *string
	CurrentStep    int
	Status         RegistrationSessionStatus
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
