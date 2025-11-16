// Package input defines input DTOs and commands for the record use cases.
package input

import "time"

// CreateRecordCommand represents input for creating a record via usecase.
type CreateRecordCommand struct {
	UserID       uint64     `json:"userId"                    validate:"required"`
	CategoryID   uint64     `json:"categoryId"                validate:"required"`
	Title        string     `json:"title"                     validate:"required,max=255"`
	Description  *string    `json:"description,omitempty"`
	TagID        uint64     `json:"tagId"                     validate:"required"`
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
	Title        *string    `json:"title,omitempty"           validate:"omitempty,max=255"`
	Description  *string    `json:"description,omitempty"`
	CategoryID   *uint64    `json:"categoryId,omitempty"`
	TagID        *uint64    `json:"tagId,omitempty"`
	EventTime    *time.Time `json:"eventTime,omitempty"`
	RecordedAt   *time.Time `json:"recordedAt,omitempty"`
	DurationSecs *int       `json:"durationSeconds,omitempty"`
	Value        *float64   `json:"value,omitempty"`
	Source       *string    `json:"source,omitempty"`
	Timezone     *string    `json:"timezone,omitempty"`
	Status       *string    `json:"status,omitempty"`
}
