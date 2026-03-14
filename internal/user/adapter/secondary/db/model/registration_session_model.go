package model

import "time"

const (
	// TableRegistrationSessions is the table name for staged registration sessions.
	TableRegistrationSessions = "aion_api.registration_sessions"
)

// RegistrationSessionDB stores staged public registration progress.
type RegistrationSessionDB struct {
	RegistrationID string    `gorm:"column:registration_id;primaryKey;type:uuid"`
	Name           string    `gorm:"column:name"`
	Username       string    `gorm:"column:username"`
	Email          string    `gorm:"column:email"`
	PasswordHash   string    `gorm:"column:password_hash"`
	Locale         *string   `gorm:"column:locale"`
	Timezone       *string   `gorm:"column:timezone"`
	Location       *string   `gorm:"column:location"`
	Bio            *string   `gorm:"column:bio"`
	AvatarURL      *string   `gorm:"column:avatar_url"`
	CurrentStep    int       `gorm:"column:current_step"`
	Status         string    `gorm:"column:status"`
	ExpiresAt      time.Time `gorm:"column:expires_at"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

// TableName overrides default table name for RegistrationSessionDB.
func (RegistrationSessionDB) TableName() string {
	return TableRegistrationSessions
}
