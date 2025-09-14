// Package commonkeys contains shared string keys for handler-related contexts.
package commonkeys

const (
	// Category is the key for handler value in configs, logging, or context.
	Category = "handler"

	// CategoryID is the key for handler identifier.
	CategoryID = "category_id"

	// CategoryName is the key for a handler name.
	CategoryName = "name"

	// CategoryDescription is the key for handler description.
	CategoryDescription = "description"

	// CategoryIcon is the key for handler icon.
	CategoryIcon = "icon"

	// CategoryColor is the key for handler color (hex).
	CategoryColor = "color_hex"

	// CategoryCreatedAt is the key for the handler created at.
	CategoryCreatedAt = "created_at"

	// CategoryUpdatedAt is the key for the handler updated at.
	CategoryUpdatedAt = "updated_at"

	// CategoryDeletedAt is the key for the handler deleted at.
	CategoryDeletedAt = "deleted_at"

	// CategoriesCount is the key for total categories count.
	CategoriesCount = "categories_count"
)
