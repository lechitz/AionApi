package category

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service provides operations for managing categories including creation, retrieval, updates, and soft deletion, utilizing a repository and logger.
type Service struct {
	Repository output.CategoryStore
	Logger     output.Logger
}

// NewCategoryService creates and returns a new instance of Service with the given repository and logger dependencies.
func NewCategoryService(repository output.CategoryStore, logger output.Logger) *Service {
	return &Service{
		Repository: repository,
		Logger:     logger,
	}
}
