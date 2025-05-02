package graph

import (
	"context"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/resolvers"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver struct used by gqlgen.
type Resolver struct {
	// Add any fields that resolvers need here
}

// NewResolver creates a new resolver instance.
func NewResolver() *Resolver {
	return &Resolver{}
}

// SymbolDefinitionResolver is for providing field resolvers on the SymbolDefinition type.
type SymbolDefinitionResolver interface {
	NextExDividendDate(ctx context.Context, obj *model.SymbolDefinition, singleFlight *bool) (*time.Time, error)
}

// SymbolDefinition returns our SymbolDefinitionResolver implementation.
func (r *Resolver) SymbolDefinition() SymbolDefinitionResolver {
	return &symbolDefinitionResolver{r}
}

type symbolDefinitionResolver struct{ *Resolver }

// NextExDividendDate delegates to our custom resolver implementation.
func (r *symbolDefinitionResolver) NextExDividendDate(ctx context.Context, obj *model.SymbolDefinition, singleFlight *bool) (*time.Time, error) {
	// Delegate to our custom implementation
	return resolvers.NextExDividendDate(ctx, obj, singleFlight)
}
