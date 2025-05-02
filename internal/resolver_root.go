package graph

import (
	"context"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/resolvers"
)

// Resolver implements the graph.ResolverRoot interface.
// It holds dependencies and acts as the receiver for resolver methods.
type Resolver struct {
	// Add dependencies here, e.g.:
	// DB *sql.DB
}

// NewResolver creates the root resolver. Called by main.go.
// It returns graph.ResolverRoot, which our *Resolver implicitly satisfies.
func NewResolver() graph.ResolverRoot {
	return &Resolver{}
}

// --- Query Resolvers ---

type queryResolver struct{ *Resolver }

// Query returns the query resolver implementation satisfying graph.QueryResolver.
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

// Symbols delegates the Query.symbols field resolution.
func (r *queryResolver) Symbols(ctx context.Context, names []string) ([]*model.SymbolDefinition, error) {
	return resolvers.SymbolsImpl(ctx, names)
}

// --- Subscription Resolvers ---

type subscriptionResolver struct{ *Resolver }

// Subscription returns the subscription resolver implementation satisfying graph.SubscriptionResolver.
func (r *Resolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}

// SymbolUpdates delegates the Subscription.symbolUpdates field resolution.
func (r *subscriptionResolver) SymbolUpdates(ctx context.Context, names []string) (<-chan *model.SymbolDefinition, error) {
	return resolvers.SymbolUpdatesImpl(ctx, names)
}

// --- SymbolDefinition Resolvers ---

type symbolDefinitionResolver struct{ *Resolver }

// SymbolDefinition returns the SymbolDefinition resolver implementation.
// The return type implicitly satisfies the generated interface.
func (r *Resolver) SymbolDefinition() graph.SymbolDefinitionResolver {
	return &symbolDefinitionResolver{r}
}

// NextExDividendDate delegates the SymbolDefinition.NextExDividendDate field resolution.
func (r *symbolDefinitionResolver) NextExDividendDate(ctx context.Context, obj *model.SymbolDefinition, singleFlight *bool) (*time.Time, error) {
	return resolvers.NextExDividendDate(ctx, obj, singleFlight)
}
