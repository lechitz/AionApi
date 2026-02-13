package cache

import "time"

// UserCacheDTO represents user data stored in cache.
type UserCacheDTO struct {
	Version   int        `json:"version"`
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Locale    *string    `json:"locale,omitempty"`
	Timezone  *string    `json:"timezone,omitempty"`
	Location  *string    `json:"location,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
