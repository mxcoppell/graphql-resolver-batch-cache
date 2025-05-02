package loaders

import (
	"context"
	"log"
	"net/http"
	"time"

	// No longer need model import here as Date maps directly to time.Time
	// model "github.com/mxcoppell/graphql-resolver-batch-cache/gen/graph/model"
	"github.com/vikstrous/dataloadgen"
)

// DividendDateLoader is a DataLoader for fetching dividend dates by symbol name
type DividendDateLoader struct {
	loader *dataloadgen.Loader[string, *time.Time]
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

// LoadDividendDate loads the dividend date for a symbol, handling singleFlight logic
func (d *DividendDateLoader) LoadDividendDate(ctx context.Context, symbolName string, singleFlight bool) (*time.Time, error) {
	if !singleFlight {
		// If singleFlight is false, check if already attempted
		if d.attemptTracker.IsAttempted(symbolName) {
			log.Printf("Symbol %s already attempted with singleFlight=false, returning nil", symbolName)
			return nil, nil // Return nil on subsequent accesses when singleFlight=false
		}
		// Mark as attempted ONLY when singleFlight is false, before the first fetch
		d.attemptTracker.MarkAttempted(symbolName)
		log.Printf("Symbol %s first attempt with singleFlight=false, marked as attempted", symbolName)
	}

	// Fetch via dataloader (handles batching and caching for singleFlight=true implicitly)
	log.Printf("Loading dividend date for %s via dataloader (singleFlight=%v)", symbolName, singleFlight)
	return d.loader.Load(ctx, symbolName)
}

// LoadManyDividendDates loads multiple dividend dates at once
// Note: This simplistic LoadMany doesn't elegantly handle the singleFlight=false logic across multiple calls.
// A more robust implementation might require modifying dataloadgen or a custom batch function.
// For now, it assumes singleFlight=true behavior for simplicity when calling LoadMany.
func (d *DividendDateLoader) LoadManyDividendDates(ctx context.Context, symbolNames []string) ([]*time.Time, []error) {
	// Mark all as attempted - This is slightly incorrect for singleFlight=false logic, as it marks before fetching.
	// This method might need revisiting depending on precise requirements for LoadMany with singleFlight=false.
	// for _, name := range symbolNames {
	// 	d.attemptTracker.MarkAttempted(name)
	// }

	// Load each one individually using the main Load method (which respects singleFlight=true implicitly)
	results := make([]*time.Time, len(symbolNames))
	errors := make([]error, len(symbolNames))

	log.Printf("LoadManyDividendDates called for %d symbols. Assuming singleFlight=true behavior.", len(symbolNames))
	for i, name := range symbolNames {
		// Assuming singleFlight=true for LoadMany for simplicity
		results[i], errors[i] = d.loader.Load(ctx, name)
	}

	return results, errors
}

// Simulate an expensive API call that fetches dividend dates for multiple symbols at once
func fetchDividendDates(ctx context.Context, symbolNames []string) ([]*time.Time, []error) {
	results := make([]*time.Time, len(symbolNames))
	errors := make([]error, len(symbolNames))

	// Simulate API latency
	time.Sleep(500 * time.Millisecond)

	// Simulate batch API response
	for i, name := range symbolNames {
		// Deterministic logic to generate dividend dates for demo purposes
		// In reality, this would call an actual API
		var date time.Time
		if name == "AAPL" {
			// One month from now
			date = time.Now().AddDate(0, 1, 0)
		} else if name == "MSFT" {
			// Two months from now
			date = time.Now().AddDate(0, 2, 0)
		} else if name == "GOOG" {
			// Three months from now
			date = time.Now().AddDate(0, 3, 0)
		} else {
			// For any other symbol, 6 months from now
			date = time.Now().AddDate(0, 6, 0)
		}
		results[i] = &date
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
