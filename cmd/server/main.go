package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	customgraph "github.com/mxcoppell/graphql-resolver-batch-cache/internal" // Alias for internal package root
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/loaders"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create a new GraphQL server using the custom resolver implementation from internal/resolver_root.go
	// Note: customgraph.NewResolver() returns graph.ResolverRoot which is compatible.
	resolver := customgraph.NewResolver() // Call the NewResolver from internal/resolver_root.go

	// Create a handler.Server manually
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Add transports (order might matter depending on routing library)
	srv.AddTransport(transport.Options{})       // Needs POST, GET, etc. - Options{} provides defaults
	srv.AddTransport(transport.GET{})           // Explicitly add GET
	srv.AddTransport(transport.POST{})          // Explicitly add POST
	srv.AddTransport(transport.MultipartForm{}) // If file uploads are needed

	// Add WebSocket support for subscriptions
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	// Enable introspection for better developer experience
	srv.Use(extension.Introspection{})

	// Create the handler chain with the dataloader middleware
	http.Handle("/", playground.Handler("GraphQL Resolver Batch Cache Demo", "/query"))
	http.Handle("/query", loaders.Middleware(srv))

	// Start the server
	log.Printf("Server running at http://localhost:%s/", port)
	log.Printf("GraphQL endpoint: http://localhost:%s/query", port)
	log.Printf("GraphQL playground: http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
