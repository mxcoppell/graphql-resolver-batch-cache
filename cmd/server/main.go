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

	// Import the graph package containing the merged resolver logic
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/graph"
	// Keep generatedGraph for the schema
	generatedGraph "github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph"
	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/loaders"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create resolver using the unified NewResolver from internal/graph/resolver.go
	resolver := graph.NewResolver() // Use the resolver from internal/graph

	// Create a handler.Server manually using the generated schema and the unified resolver
	srv := handler.New(generatedGraph.NewExecutableSchema(generatedGraph.Config{Resolvers: resolver}))

	// Add transports (order might matter depending on routing library)
	srv.AddTransport(transport.Options{})       // Needs POST, GET, etc. - Options{} provides defaults
	srv.AddTransport(transport.GET{})           // Explicitly add GET
	srv.AddTransport(transport.POST{})          // Explicitly add POST
	srv.AddTransport(transport.MultipartForm{}) // If file uploads are needed

	// Add WebSocket support for subscriptions
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for simplicity in demo
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
