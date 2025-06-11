// Package security provides utility functions for JWT token generation and validation.
package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserID is a constant key used to represent a user's unique identifier in various operations, such as token generation or data storage.
const UserID = "user_id"

// Exp is a constant key that represents the expiration claim in a JWT (JSON Web Token).
const Exp = "exp"

// ExpTimeToken defines the duration of 24 hours used as the standard token expiration period in time-based operations.
const ExpTimeToken = 24 * time.Hour

// GenerateToken creates a signed JWT token with userID and expiration using the provided secretKey.
func GenerateToken(userID uint64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		UserID: userID,
		Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
