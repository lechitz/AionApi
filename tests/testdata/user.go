package testdata

import (
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// TestPerfectUser is a predefined UserDomain instance used for testing, representing a fully initialized user with typical attributes and lifecycle timestamps.
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
