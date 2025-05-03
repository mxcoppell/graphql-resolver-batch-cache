package graph

import (
	"context"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/resolvers"
)

// symbolDefinitionResolver implements the generated SymbolDefinitionResolver interface.
type symbolDefinitionResolver struct{ *Resolver }

// NextExDividendDate delegates to our custom resolver implementation.
func (r *symbolDefinitionResolver) NextExDividendDate(ctx context.Context, obj *model.SymbolDefinition, singleFlight *bool) (*time.Time, error) {
	// Delegate to our custom implementation
	return resolvers.NextExDividendDate(ctx, obj, singleFlight)
}
