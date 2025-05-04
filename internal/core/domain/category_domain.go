package domain

import "time"

type Category struct {
	ID          uint64
	UserID      uint64
	Name        string
	Description string
	Color       string
	Icon        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
