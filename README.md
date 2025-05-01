# GraphQL Resolver Batch Cache Demo

This project demonstrates a pattern for efficiently managing cached resolver data in GraphQL, specifically for the case of a resolver field `NextExDividendDate` in a `SymbolDefinition` type. It shows how to:

1. Use a DataLoader to batch and cache requests for dividend dates
2. Allow clients to control caching behavior with a `useCache` parameter
3. Handle both Query and Subscription operations

## Key Features

- **DataLoader Pattern** - Consolidates multiple resolver calls into a single batch request, avoiding the N+1 query problem
- **Request-Scoped Caching** - The DataLoader caches results within a single request/subscription cycle
- **Attempt Tracking** - Tracks which symbols have been resolved, enabling clients to control whether subsequent accesses should return cached values or nil

## Resolver Behavior

The `NextExDividendDate` field resolver has the following behavior:

- When `useCache` is `true` (default):
  - First access: Makes the upstream API call and caches the result
  - Subsequent accesses: Returns the cached value

- When `useCache` is `false`:
  - First access: Makes the upstream API call and marks the symbol as attempted
  - Subsequent accesses: Returns `nil`

This is particularly useful in subscription scenarios where you might want a field to be resolved once but return `nil` on subsequent resolver invocations during the same subscription event.

## Getting Started

### Prerequisites

- Go 1.23.x or later

### Installation

1. Clone the repository
2. Install dependencies:

```bash
go mod tidy
go mod vendor
```

### Running the Server

```bash
go run cmd/server/main.go
```

The server will start at http://localhost:8080 with a GraphQL playground.

## Example Queries

### Basic Query

```graphql
query GetSymbols {
  symbols(names: ["AAPL", "MSFT", "GOOG"]) {
    Name
    NextExDividendDate
  }
}
```

### Query with useCache Parameter

```graphql
query GetSymbols {
  symbols(names: ["AAPL", "MSFT", "GOOG"]) {
    Name
    NextExDividendDate(useCache: false)
  }
}
```

### Subscription

```graphql
subscription WatchSymbols {
  symbolUpdates(names: ["AAPL", "MSFT", "GOOG"]) {
    Name
    NextExDividendDate(useCache: false)
  }
}
```

## Implementation Details

- The DataLoader is defined in `graph/loaders/dataloaders.go`
- The HTTP middleware to inject the loader into each request is in the same file
- The resolver logic is in `graph/schema.resolvers.go`
- The server configuration is in `cmd/server/main.go`

## License

MIT 