package resolvers

import (
	"context"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/loaders"
)

// NextExDividendDate resolves the NextExDividendDate field for the SymbolDefinition type.
// This custom implementation is used by registering it with the ResolverRoot.
// It replaces the auto-generated resolver.
func NextExDividendDate(ctx context.Context, obj *model.SymbolDefinition, singleFlight *bool) (*time.Time, error) {
	// Get the loader from context
	loader := loaders.For(ctx)

	// Determine the singleFlight flag value (default to true if not specified)
	shouldSingleFlight := true
	if singleFlight != nil {
		shouldSingleFlight = *singleFlight
	}

	// Load the dividend date using the loader, passing the singleFlight flag
	dateResult, err := loader.LoadDividendDate(ctx, obj.Name, shouldSingleFlight)
	if err != nil {
		return nil, err
	}

	return dateResult, nil // Loader now returns *time.Time directly
}
