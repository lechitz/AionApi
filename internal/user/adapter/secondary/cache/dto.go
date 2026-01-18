package cache

import "time"

// UserCacheDTO represents user data stored in cache.
// SECURITY: Password hash is NEVER cached - always fetch from database for authentication.
type UserCacheDTO struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// ⛔ NEVER CACHE: PasswordHash
	// Password verification ALWAYS goes to database
}
