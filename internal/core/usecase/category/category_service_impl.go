package category

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Service provides operations for managing categories including creation, retrieval, updates, and soft deletion, utilizing a repository and logger.
type Service struct {
	Repository db.CategoryStore
	Logger     logger.Logger
}

// NewCategoryService creates and returns a new instance of Service with the given repository and logger dependencies.
func NewCategoryService(repository db.CategoryStore, logger logger.Logger) *Service {
	return &Service{
		Repository: repository,
		Logger:     logger,
	}
}
