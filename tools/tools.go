//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/google/uuid" // Add other tool dependencies if needed
	_ "github.com/vikstrous/dataloadgen"
)
