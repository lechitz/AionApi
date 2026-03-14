package output

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/user/core/domain"
)

// RegistrationSessionUniqueness holds uniqueness checks across final users and pending sessions.
type RegistrationSessionUniqueness struct {
	UsernameTaken bool
	EmailTaken    bool
}

// RegistrationSessionRepository defines persistence operations for staged public registrations.
type RegistrationSessionRepository interface {
	CreateRegistrationSession(ctx context.Context, session domain.RegistrationSession) (domain.RegistrationSession, error)
	GetRegistrationSessionByID(ctx context.Context, registrationID string) (domain.RegistrationSession, error)
	UpdateRegistrationSession(ctx context.Context, registrationID string, fields map[string]interface{}) (domain.RegistrationSession, error)
	DeleteRegistrationSession(ctx context.Context, registrationID string) error
	CheckRegistrationUniqueness(ctx context.Context, username, email string, now time.Time) (RegistrationSessionUniqueness, error)
}
