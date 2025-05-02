package resolvers

import (
	"context"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
)

// Context is an alias for the standard context.Context
type Context = context.Context

// SymbolDefinition is an alias for the generated model type
type SymbolDefinition = *model.SymbolDefinition
