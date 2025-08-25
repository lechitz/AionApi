// Package category implements the GraphQL controller (handler) for Category operations.
// It orchestrates tracing, structured logging, and mapping, then delegates business logic
// to the core use case via the input port (CategoryService).
package category

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Handler is responsible for handling GraphQL Category mutations and queries.
// It wraps requests in tracing spans, executes mapping and logging logic,
// and calls the core business logic defined in CategoryService.
type Handler struct {
	CategoryService input.CategoryService
	Logger          output.ContextLogger
}

// NewHandler constructs a new Handler with the provided CategoryService and Logger.
func NewHandler(categoryService input.CategoryService, logger output.ContextLogger) *Handler {
	return &Handler{CategoryService: categoryService, Logger: logger}
}
