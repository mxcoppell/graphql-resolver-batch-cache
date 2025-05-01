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
	"github.com/mxcoppell/graphql-resolver-batch-cache/graph"
	"github.com/mxcoppell/graphql-resolver-batch-cache/graph/loaders"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create a new GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

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
