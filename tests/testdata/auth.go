// Package testdata contains test data used for testing purposes.
package testdata

import (
	"github.com/lechitz/AionApi/internal/core/domain/entity"
)

// TestPerfectLoginInputUser is a predefined instance of UserDomain representing a user with valid login credentials for testing purposes.
var TestPerfectLoginInputUser = entity.UserDomain{
	Username: "testuser",
	Password: "password123",
}

// HashedPassword stores a hashed representation of a password using bcrypt or similar hashing algorithms.
var HashedPassword = "$2a$12$hashedExample..."
