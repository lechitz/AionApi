// Package input defines input DTOs and commands for the record use cases.
package input

import "time"

// CreateRecordCommand represents input for creating a record via usecase.
// Note: category is obtained via Tag relationship (Record → Tag → Category).
type CreateRecordCommand struct {
	UserID       uint64     `json:"userId"                    validate:"required"`
	TagID        uint64     `json:"tagId"                     validate:"required"`
	Description  *string    `json:"description,omitempty"`
	EventTime    time.Time  `json:"eventTime"                 validate:"required"`
	RecordedAt   *time.Time `json:"recordedAt,omitempty"`
	DurationSecs *int       `json:"durationSeconds,omitempty"`
	Value        *float64   `json:"value,omitempty"`
	Source       *string    `json:"source,omitempty"`
	Timezone     *string    `json:"timezone,omitempty"`
	Status       *string    `json:"status,omitempty"`
}

// UpdateRecordCommand represents fields allowed to be updated.
type UpdateRecordCommand struct {
	Description  *string    `json:"description,omitempty"`
	TagID        *uint64    `json:"tagId,omitempty"`
	EventTime    *time.Time `json:"eventTime,omitempty"`
	RecordedAt   *time.Time `json:"recordedAt,omitempty"`
	DurationSecs *int       `json:"durationSeconds,omitempty"`
	Value        *float64   `json:"value,omitempty"`
	Source       *string    `json:"source,omitempty"`
	Timezone     *string    `json:"timezone,omitempty"`
	Status       *string    `json:"status,omitempty"`
}
