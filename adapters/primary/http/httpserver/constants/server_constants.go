// Package constants contains constants used throughout the application.
package constants

// ErrContextPathEmpty is returned when the context path is empty.
const ErrContextPathEmpty = "contextPath cannot be empty"

// ErrContextPathSlash is returned when the context path contains additional slashes (`/`).
const ErrContextPathSlash = "contextPath cannot contain additional slashes `/`"

// InvalidContextPath is returned when the context path is invalid. It uses formatting for path insertion.
const InvalidContextPath = "invalid context path: '%s'"
