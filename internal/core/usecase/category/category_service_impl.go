package category

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Service provides operations for managing categories including creation, retrieval, updates, and soft deletion, utilizing a repository and logger.
type Service struct {
	Repository output.CategoryStore
	Logger     logger.Logger
}

// NewCategoryService creates and returns a new instance of Service with the given repository and logger dependencies.
func NewCategoryService(repository output.CategoryStore, logger logger.Logger) *Service {
	return &Service{
		Repository: repository,
		Logger:     logger,
	}
}
