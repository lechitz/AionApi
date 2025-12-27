// Package output defines interfaces for tag-related cache operations.
package output

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/tag/core/domain"
)

// TagCache defines cache operations for tags.
type TagCache interface {
	SaveTag(ctx context.Context, tag domain.Tag, expiration time.Duration) error
	SaveTagByName(ctx context.Context, tag domain.Tag, expiration time.Duration) error
	SaveTagList(ctx context.Context, userID uint64, tags []domain.Tag, expiration time.Duration) error
	SaveTagsByCategory(ctx context.Context, categoryID, userID uint64, tags []domain.Tag, expiration time.Duration) error
	GetTag(ctx context.Context, tagID, userID uint64) (domain.Tag, error)
	GetTagByName(ctx context.Context, tagName string, userID uint64) (domain.Tag, error)
	GetTagList(ctx context.Context, userID uint64) ([]domain.Tag, error)
	GetTagsByCategory(ctx context.Context, categoryID, userID uint64) ([]domain.Tag, error)
	DeleteTag(ctx context.Context, tagID, userID uint64) error
	DeleteTagByName(ctx context.Context, tagName string, userID uint64) error
	DeleteTagList(ctx context.Context, userID uint64) error
	DeleteTagsByCategory(ctx context.Context, categoryID, userID uint64) error
}
