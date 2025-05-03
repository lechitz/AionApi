package graph

import (
	"github.com/lechitz/AionApi/internal/core/ports/input/graphql"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

type Resolver struct {
	CategoryService graphql.CategoryService
	Logger          logger.Logger
}
