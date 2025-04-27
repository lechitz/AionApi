package domain

import "time"

type CategoryDomain struct {
	ID          int64
	Name        string
	Description string
	Color       string
	Icon        string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
