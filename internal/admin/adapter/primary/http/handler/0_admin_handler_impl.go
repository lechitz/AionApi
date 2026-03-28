// Package handler is the handler for the admin context in the application.
package handler

import (
	"github.com/lechitz/aion-api/internal/admin/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// Handler is the handler for admin-related HTTP operations.
type Handler struct {
	AdminService input.AdminService
	Logger       logger.ContextLogger
	Config       *config.Config
}

// New returns an Admin handler with dependencies injected.
func New(adminService input.AdminService, cfg *config.Config, logger logger.ContextLogger) *Handler {
	return &Handler{
		AdminService: adminService,
		Config:       cfg,
		Logger:       logger,
	}
}
