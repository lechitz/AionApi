package input

// CreateTagCommand represents the data required to create a new tag.
type CreateTagCommand struct {
	Name        string
	ID          uint64
	UserID      uint64
	CategoryID  uint64
	Description *string
	Icon        *string
}

// UpdateTagCommand represents the data required to update an existing tag.
type UpdateTagCommand struct {
	Name        *string
	Description *string
	CategoryID  *uint64
	ID          uint64
	UserID      uint64
	Icon        *string
}
