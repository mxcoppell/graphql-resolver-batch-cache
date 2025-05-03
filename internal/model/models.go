package model

import "time"

// SymbolDefinition represents a financial instrument symbol
type SymbolDefinition struct {
	Name string `json:"name"`
}

// Date is a wrapper around time.Time that implements the graphql.Marshaler interface
type Date struct {
	time.Time
}
