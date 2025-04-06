package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"gorm.io/gorm"
	"time"
)

var TestPerfectAuthUser = domain.UserDomain{
	ID:        1,
	Name:      "Auth Test User",
	Username:  "user",
	Email:     "user@example.com",
	Password:  "supersecure123",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	DeletedAt: gorm.DeletedAt{},
}
