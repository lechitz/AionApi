package usecase

import (
	"github.com/lechitz/AionApi/internal/category/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Service provides operations for managing categories including creation, retrieval, updates, and soft deletion, using a repository and contextlogger.
type Service struct {
	CategoryRepository output.CategoryRepository
	Logger             logger.ContextLogger
}

// NewService creates and returns a new instance of Service with the given repository and contextlogger dependencies.
func NewService(repository output.CategoryRepository, logger logger.ContextLogger) *Service {
	return &Service{
		CategoryRepository: repository,
		Logger:             logger,
	}
}
