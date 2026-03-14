// Package usecase contains business logic for the audit context.
package usecase

import (
	"github.com/lechitz/AionApi/internal/audit/core/ports/input"
	"github.com/lechitz/AionApi/internal/audit/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Service provides audit event write/list operations.
type Service struct {
	repository output.AuditActionEventRepository
	logger     logger.ContextLogger
}

// NewService creates a new audit service implementation.
func NewService(repository output.AuditActionEventRepository, log logger.ContextLogger) input.Service {
	return &Service{
		repository: repository,
		logger:     log,
	}
}
