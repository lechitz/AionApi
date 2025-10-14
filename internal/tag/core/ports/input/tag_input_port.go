// Package input is the input port for the tag context.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/tag/core/domain"
)

// TagCreator interface for create a new tag.
type TagCreator interface {
	Create(ctx context.Context, tag CreateTagCommand) (domain.Tag, error)
}

// TagRetriever defines methods for retrieving tag details or multiple tags for a specific user.
type TagRetriever interface {
	GetByID(ctx context.Context, tagID uint64, userID uint64) (domain.Tag, error)
	GetByName(ctx context.Context, tagName string, userID uint64) (domain.Tag, error)
}

// TagService defines a contract that combines creating, retrieving, updating, and soft-deleting tags in the system.
type TagService interface {
	TagCreator
	TagRetriever
}
