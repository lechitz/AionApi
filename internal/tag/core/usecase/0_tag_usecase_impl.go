// Package usecase contains the business logic for the Tag context.
// It orchestrates input commands, applies validation and domain rules,
// and delegates persistence to the repository layer while handling
// observability and logging concerns.package usecase
package usecase

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/tag/core/ports/output"
)

// Service provides operations for managing tags including creation, retrieval, updates, and soft deletion, using a repository and contextlogger.
type Service struct {
	TagRepository output.TagRepository
	TagCache      output.TagCache
	Logger        logger.ContextLogger
}

// NewService creates and returns a new instance of Service with the given repository and contextlogger dependencies.
func NewService(tagRepository output.TagRepository, tagCache output.TagCache, logger logger.ContextLogger) *Service {
	return &Service{
		TagRepository: tagRepository,
		TagCache:      tagCache,
		Logger:        logger,
	}
}
