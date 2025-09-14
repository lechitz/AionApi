package input

// CreateUserCommand defines a struct for creating a new user.
type CreateUserCommand struct {
	Name     string
	Username string
	Email    string
	Password string
}

// UpdateUserCommand defines a struct for updating user information.
type UpdateUserCommand struct {
	Name     *string
	Username *string
	Email    *string
}

// HasUpdates returns true if the command contains any updates.
func (u UpdateUserCommand) HasUpdates() bool {
	return u.Name != nil || u.Username != nil || u.Email != nil
}
