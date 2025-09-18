package input

// CreateTagCommand represents the data required to create a new tag.
type CreateTagCommand struct {
	Name        string
	ID          uint64
	UserID      uint64
	CategoryID  uint64
	Description *string
}
