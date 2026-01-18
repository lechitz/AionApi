// Package usecase (admin) provides operations for managing admin functions.
package usecase

import (
	"github.com/lechitz/AionApi/internal/admin/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Service provides an abstraction for admin management operations.
type Service struct {
	adminRepository output.AdminRepository
	logger          logger.ContextLogger
}

// NewService creates and returns a new Service instance with the provided dependencies.
func NewService(
	adminRepository output.AdminRepository,
	logger logger.ContextLogger,
) *Service {
	return &Service{
		adminRepository: adminRepository,
		logger:          logger,
	}
}
