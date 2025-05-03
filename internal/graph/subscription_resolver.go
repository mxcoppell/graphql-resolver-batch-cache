package graph

import (
	"context"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/resolvers"
)

// subscriptionResolver implements the generated SubscriptionResolver interface.
type subscriptionResolver struct{ *Resolver }

// SymbolUpdates delegates the Subscription.symbolUpdates field resolution.
func (r *subscriptionResolver) SymbolUpdates(ctx context.Context, names []string) (<-chan *model.SymbolDefinition, error) {
	return resolvers.SymbolUpdatesImpl(ctx, names)
}
