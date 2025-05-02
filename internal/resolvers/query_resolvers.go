package resolvers

import (
	"context"
	"fmt"
	"log"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
)

// SymbolsImpl provides the implementation logic for the Query.symbols resolver.
func SymbolsImpl(ctx context.Context, names []string) ([]*model.SymbolDefinition, error) {
	log.Printf("Query.symbols called with %d symbols", len(names))

	// Create symbol definitions for each name
	result := make([]*model.SymbolDefinition, len(names))
	for i, name := range names {
		// Here we create model objects with just the name filled in
		// NextExDividendDate will be resolved separately when requested
		desc := fmt.Sprintf("Description for Symbol %s", name)
		result[i] = &model.SymbolDefinition{
			Name:        name,
			Description: &desc,
		}
	}

	return result, nil
}
