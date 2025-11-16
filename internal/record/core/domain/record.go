// Package domain contains domain entities for the record bounded context.
package domain

import "time"

// Record represents a single logged event in the system (one-off or transformed from planner/habit).
type Record struct {
	ID          uint64  `json:"id"                    db:"id"`
	UserID      uint64  `json:"userId"                db:"user_id"`
	Title       string  `json:"title"                 db:"title"`
	Description *string `json:"description,omitempty" db:"description"`
	CategoryID  uint64  `json:"categoryId"            db:"category_id"`

	EventTime time.Time `json:"eventTime" db:"event_time"` // when the event was planned/scheduled to occur

	RecordedAt *time.Time `json:"recordedAt,omitempty" db:"recorded_at"` // when the event was actually logged/recorded

	DurationSecs *int     `json:"durationSeconds,omitempty" db:"duration_seconds"`
	Value        *float64 `json:"value,omitempty"           db:"value"`
	Source       *string  `json:"source,omitempty"          db:"source"`
	Timezone     *string  `json:"timezone,omitempty"        db:"timezone"`
	Status       *string  `json:"status,omitempty"          db:"status"`

	CreatedAt time.Time  `json:"createdAt"           db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt"           db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

	TagID uint64 `json:"tagId" db:"tag_id"`
}
