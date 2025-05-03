package graph

// Use the generated graph package for the interface types
import (
	generatedGraph "github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph"
	// We don't need context, time, model, or resolvers here anymore
	// as the specific implementations are moved to other files.
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver struct used by gqlgen.
// It implements the generatedGraph.ResolverRoot interface.
type Resolver struct {
	// Add any fields that resolvers need here, e.g., database connections
}

// NewResolver creates a new resolver instance.
func NewResolver() generatedGraph.ResolverRoot { // Return the generated interface type
	return &Resolver{}
}

// Query returns the query resolver implementation satisfying generatedGraph.QueryResolver.
func (r *Resolver) Query() generatedGraph.QueryResolver { // Use generated interface type
	return &queryResolver{r}
}

// Subscription returns the subscription resolver implementation satisfying generatedGraph.SubscriptionResolver.
func (r *Resolver) Subscription() generatedGraph.SubscriptionResolver { // Use generated interface type
	return &subscriptionResolver{r}
}

// SymbolDefinition returns our SymbolDefinitionResolver implementation.
// It satisfies the generatedGraph.SymbolDefinitionResolver interface.
func (r *Resolver) SymbolDefinition() generatedGraph.SymbolDefinitionResolver { // Use generated interface type
	return &symbolDefinitionResolver{r}
}
