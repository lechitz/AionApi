package constants

import "time"

// JWTKeyLength defines the length of the JWT key.
const JWTKeyLength = 64

// ExpTimeToken defines the duration of 24 hours used as the standard token expiration period in time-based operations.
const ExpTimeToken = 24 * time.Hour
