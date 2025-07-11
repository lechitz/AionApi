// Package graph The package graph provides the GraphQL resolver implementation.
//
//nolint:govet,revive, perfsprint,nolintlint
package graph

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

type Resolver struct {
	CategoryService input.CategoryService
	Logger          output.ContextLogger
}
