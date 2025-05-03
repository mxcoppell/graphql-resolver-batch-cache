package graph

import (
	"context"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/resolvers"
)

// queryResolver implements the generated QueryResolver interface.
type queryResolver struct{ *Resolver }

// Symbols delegates the Query.symbols field resolution.
func (r *queryResolver) Symbols(ctx context.Context, names []string) ([]*model.SymbolDefinition, error) {
	return resolvers.SymbolsImpl(ctx, names)
}
