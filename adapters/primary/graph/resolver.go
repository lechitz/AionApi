package graph

import (
	"github.com/lechitz/AionApi/internal/core/ports/input/graphql"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Resolver struct serves as the dependency injection container for GraphQL resolvers.
// It provides access to category services and logging functionalities.
type Resolver struct {
	CategoryService graphql.CategoryService
	Logger          logger.Logger
}
