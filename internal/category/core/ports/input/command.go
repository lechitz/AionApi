// Package input defines the input commands (DTOs) used to interact with Category use cases.
package input

// CreateCategoryCommand represents the data required to create a new category.
type CreateCategoryCommand struct {
	Description *string
	ColorHex    *string
	Icon        *string
	Name        string
	UserID      uint64
}

// UpdateCategoryCommand represents the data required to update an existing category.
type UpdateCategoryCommand struct {
	Name        *string
	Description *string
	ColorHex    *string
	Icon        *string
	ID          uint64
	UserID      uint64
}
