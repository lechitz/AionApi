package testdata

import (
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
)

var TestPerfectUser = domain.UserDomain{
	ID:        1,
	Name:      "User Name",
	Username:  "user",
	Email:     "user@example.com",
	Password:  "supersecure123",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	DeletedAt: nil,
}
