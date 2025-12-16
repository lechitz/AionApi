package input

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

// UpdateUserCommand defines a struct for updating user information.
type UpdateUserCommand struct {
	Name      *string
	Username  *string
	Email     *string
	Locale    *string
	Timezone  *string
	Location  *string
	Bio       *string
	AvatarURL *string
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
		u.AvatarURL != nil
}
