package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/graph/model"
	"github.com/vikstrous/dataloadgen"
)

// DividendDateLoader is a DataLoader for fetching dividend dates by symbol name
type DividendDateLoader struct {
	loader *dataloadgen.Loader[string, *model.Date]
	// Track which symbols have already been attempted in this request/subscription cycle
	attemptTracker *SymbolAttemptTracker
}

// SymbolAttemptTracker tracks which symbol names have been attempted in this request
type SymbolAttemptTracker struct {
	attemptedSymbols map[string]bool
}

// NewSymbolAttemptTracker creates a new tracker
func NewSymbolAttemptTracker() *SymbolAttemptTracker {
	return &SymbolAttemptTracker{
		attemptedSymbols: make(map[string]bool),
	}
}

// IsAttempted checks if a symbol has been attempted
func (t *SymbolAttemptTracker) IsAttempted(symbol string) bool {
	return t.attemptedSymbols[symbol]
}

// MarkAttempted marks a symbol as attempted
func (t *SymbolAttemptTracker) MarkAttempted(symbol string) {
	t.attemptedSymbols[symbol] = true
}

// NewDividendDateLoader creates a new DividendDateLoader
func NewDividendDateLoader() *DividendDateLoader {
	return &DividendDateLoader{
		loader:         dataloadgen.NewLoader(fetchDividendDates),
		attemptTracker: NewSymbolAttemptTracker(),
	}
}

// LoadDividendDate loads the dividend date for a symbol, with the option to respect attempted tracking
func (d *DividendDateLoader) LoadDividendDate(ctx context.Context, symbolName string, useCache bool) (*model.Date, error) {
	// If useCache is false and we've already attempted this symbol, return nil
	if !useCache && d.attemptTracker.IsAttempted(symbolName) {
		return nil, nil
	}

	// Mark as attempted regardless of outcome
	d.attemptTracker.MarkAttempted(symbolName)

	// Fetch via dataloader (which will batch and cache)
	return d.loader.Load(ctx, symbolName)
}

// LoadManyDividendDates loads multiple dividend dates at once
func (d *DividendDateLoader) LoadManyDividendDates(ctx context.Context, symbolNames []string) ([]*model.Date, []error) {
	// Mark all as attempted
	for _, name := range symbolNames {
		d.attemptTracker.MarkAttempted(name)
	}

	// Load each one individually since dataloadgen doesn't have a LoadMany method
	results := make([]*model.Date, len(symbolNames))
	errors := make([]error, len(symbolNames))

	for i, name := range symbolNames {
		results[i], errors[i] = d.loader.Load(ctx, name)
	}

	return results, errors
}

// Simulate an expensive API call that fetches dividend dates for multiple symbols at once
func fetchDividendDates(ctx context.Context, symbolNames []string) ([]*model.Date, []error) {
	results := make([]*model.Date, len(symbolNames))
	errors := make([]error, len(symbolNames))

	// Simulate API latency
	time.Sleep(500 * time.Millisecond)

	// Simulate batch API response
	for i, name := range symbolNames {
		// Deterministic logic to generate dividend dates for demo purposes
		// In reality, this would call an actual API
		if name == "AAPL" {
			// One month from now
			date := time.Now().AddDate(0, 1, 0)
			results[i] = &model.Date{Time: date}
		} else if name == "MSFT" {
			// Two months from now
			date := time.Now().AddDate(0, 2, 0)
			results[i] = &model.Date{Time: date}
		} else if name == "GOOG" {
			// Three months from now
			date := time.Now().AddDate(0, 3, 0)
			results[i] = &model.Date{Time: date}
		} else {
			// For any other symbol, 6 months from now
			date := time.Now().AddDate(0, 6, 0)
			results[i] = &model.Date{Time: date}
		}
	}

	return results, errors
}

// Context key for the loader
type contextKey string

// LoaderKey is the key for the loader in the context
const LoaderKey = contextKey("dividendDateLoader")

// Middleware adds the dataloader to the context
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a loader for this request
		loader := NewDividendDateLoader()

		// Add it to the context
		ctx := context.WithValue(r.Context(), LoaderKey, loader)

		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// For returns the loader from the context
func For(ctx context.Context) *DividendDateLoader {
	return ctx.Value(LoaderKey).(*DividendDateLoader)
}
