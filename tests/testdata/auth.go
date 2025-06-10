package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

// TestPerfectLoginInputUser is a predefined instance of UserDomain representing a user with valid login credentials for testing purposes.
var TestPerfectLoginInputUser = domain.UserDomain{
	Username: "testuser",
	Password: "password123",
}

// HashedPassword stores a hashed representation of a password using bcrypt or similar hashing algorithms.
var HashedPassword = "$2a$12$hashedExample..."
