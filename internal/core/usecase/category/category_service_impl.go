package category

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Service provides operations for managing categories including creation, retrieval, updates, and soft deletion, using a repository and contextlogger.
type Service struct {
	CategoryRepository output.CategoryStore
	Logger             output.ContextLogger
}

// NewService creates and returns a new instance of Service with the given repository and contextlogger dependencies.
func NewService(repository output.CategoryStore, logger output.ContextLogger) *Service {
	return &Service{
		CategoryRepository: repository,
		Logger:             logger,
	}
}
