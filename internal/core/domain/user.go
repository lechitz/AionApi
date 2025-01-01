package domain

import "time"

type UserDomain struct {
	ID        uint64
	Name      string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
