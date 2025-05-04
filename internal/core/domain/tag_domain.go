package domain

import "time"

type Tag struct {
	ID          uint64
	UserID      uint64
	CategoryID  uint64
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
