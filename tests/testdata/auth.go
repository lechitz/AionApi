package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

var TestPerfectLoginInputUser = domain.UserDomain{
	Username: "testuser",
	Password: "password123",
}

var HashedPassword = "$2a$12$hashedExample..."
