package input

import "mime/multipart"

// CreateUserCommand defines a struct for creating a new user.
type CreateUserCommand struct {
	Name      string
	Username  string
	Email     string
	Password  string
	Locale    *string
	Timezone  *string
	Location  *string
	Bio       *string
	AvatarURL *string
}

// StartRegistrationCommand defines input data for registration step 1.
type StartRegistrationCommand struct {
	Name     string
	Username string
	Email    string
	Password string
}

// UpdateRegistrationProfileCommand defines input data for registration step 2.
type UpdateRegistrationProfileCommand struct {
	Locale   string
	Timezone string
	Location string
	Bio      string
}

// UpdateRegistrationAvatarCommand defines input data for registration step 3.
type UpdateRegistrationAvatarCommand struct {
	AvatarURL *string
}

// UploadAvatarCommand defines input data for uploading a user avatar image.
type UploadAvatarCommand struct {
	File        multipart.File
	Filename    string
	SizeBytes   int64
	ContentType string
}

// UpdateUserCommand defines a struct for updating user information.
type UpdateUserCommand struct {
	Name                *string
	Username            *string
	Email               *string
	Locale              *string
	Timezone            *string
	Location            *string
	Bio                 *string
	AvatarURL           *string
	OnboardingCompleted *bool
}

// HasUpdates returns true if the command contains any updates.
func (u UpdateUserCommand) HasUpdates() bool {
	return u.Name != nil ||
		u.Username != nil ||
		u.Email != nil ||
		u.Locale != nil ||
		u.Timezone != nil ||
		u.Location != nil ||
		u.Bio != nil ||
		u.AvatarURL != nil ||
		u.OnboardingCompleted != nil
}
