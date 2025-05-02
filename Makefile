# Makefile for the GraphQL server project

.PHONY: gen build clean

# Target to regenerate GraphQL code using gqlgen
gen:
	go run github.com/99designs/gqlgen generate

# Target to build the Go project
# Assumes the main package is in cmd/server/main.go and the output binary is named 'server'
build:
	go build -o ./bin/server ./cmd/server/main.go 

# Target to clean up generated files and build artifacts
clean:
	@echo "Cleaning up generated files and build artifacts..."
	@rm -rf ./internal/gen
	@rm -f ./bin/server
	@echo "Cleanup complete." 