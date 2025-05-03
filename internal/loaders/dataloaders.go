package loaders

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/cache"
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
	// Check if the symbol has already been attempted *in this request scope*
	alreadyAttempted := d.attemptTracker.IsAttempted(symbolName)

	// Early exit ONLY if singleFlight=true AND it was already attempted.
	if singleFlight && alreadyAttempted {
		log.Printf("Symbol %s already attempted in this scope with singleFlight=true, returning nil", symbolName)
		return nil, nil
	}

	// Mark as attempted on the first encounter within this request scope.
	// This ensures subsequent singleFlight=true calls for the same key return nil.
	if !alreadyAttempted {
		d.attemptTracker.MarkAttempted(symbolName)
		log.Printf("Symbol %s first attempt in this scope (singleFlight=%t), marked. Proceeding to dataloader.", symbolName, singleFlight)
	} else {
		// Log if it was already attempted but singleFlight is false (will proceed to dataloader)
		log.Printf("Symbol %s already attempted in this scope, but singleFlight=false. Proceeding to dataloader.", symbolName)
	}

	// Proceed to the dataloader.
	// - If first attempt: dataloader might miss, triggering batch function (which checks shared cache).
	// - If already attempted & singleFlight=false: dataloader should hit its internal request-scoped cache.
	return d.loader.Load(ctx, symbolName)
}

// LoadManyDividendDates loads multiple dividend dates at once
// Note: This simplistic LoadMany doesn't elegantly handle the singleFlight=false logic across multiple calls.
// A more robust implementation might require modifying dataloadgen or a custom batch function.
// For now, it assumes singleFlight=true behavior for simplicity when calling LoadMany.
func (d *DividendDateLoader) LoadManyDividendDates(ctx context.Context, symbolNames []string) ([]*time.Time, []error) {
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

// fetchDividendDates is the batch function used by dataloadgen.
// It now checks a shared cache before simulating the API call.
func fetchDividendDates(ctx context.Context, symbolNames []string) ([]*time.Time, []error) {
	log.Printf("DataLoader Batch Function called for keys: %v", symbolNames)
	results := make([]*time.Time, len(symbolNames))
	errors := make([]error, len(symbolNames)) // Initialize error slice

	// --- Check Shared Cache First ---
	keysToFetchFromApi := make([]string, 0, len(symbolNames))
	// Map API fetch index back to original results index
	apiFetchIndexToOrigIndex := make(map[int]int, len(symbolNames))

	for i, name := range symbolNames {
		if cachedVal, found := cache.Get(name); found {
			log.Printf("Shared cache HIT for key: %s", name)
			results[i] = cachedVal
		} else {
			log.Printf("Shared cache MISS for key: %s", name)
			apiFetchIndexToOrigIndex[len(keysToFetchFromApi)] = i
			keysToFetchFromApi = append(keysToFetchFromApi, name)
		}
	}

	// --- Fetch Missing Keys from Simulated API ---
	if len(keysToFetchFromApi) > 0 {
		log.Printf("Calling simulated API for keys: %v", keysToFetchFromApi)

		// Simulate API latency only if we need to fetch
		time.Sleep(500 * time.Millisecond)

		// Simulate batch API response for the missing keys
		apiResults := make([]*time.Time, len(keysToFetchFromApi))
		// Simulate potential API errors (can be nil)
		apiErrors := make([]error, len(keysToFetchFromApi))

		for i, name := range keysToFetchFromApi {
			// Deterministic logic for demo purposes
			var date time.Time
			if name == "AAPL" {
				log.Printf("Simulating API fetch for %s", name)
				date = time.Now().AddDate(0, 1, 0)
			} else if name == "MSFT" {
				log.Printf("Simulating API fetch for %s", name)
				date = time.Now().AddDate(0, 2, 0)
			} else if name == "GOOG" {
				log.Printf("Simulating API fetch for %s", name)
				date = time.Now().AddDate(0, 3, 0)
			} else {
				log.Printf("Simulating API fetch for %s", name)
				date = time.Now().AddDate(0, 6, 0)
			}
			// Only store non-nil results in the results slice for the dataloader
			if apiErrors[i] == nil {
				apiResults[i] = &date
				// Add successful results to the shared cache
				log.Printf("Adding API result for %s to shared cache", name)
				cache.Set(name, &date)
			} else {
				log.Printf("Simulated API error for %s: %v", name, apiErrors[i])
			}
		}

		// --- Populate main results slice from API results ---
		for apiIdx, apiRes := range apiResults {
			origIdx := apiFetchIndexToOrigIndex[apiIdx]
			if apiErrors[apiIdx] != nil {
				errors[origIdx] = apiErrors[apiIdx]
			} else {
				results[origIdx] = apiRes
			}
		}
	}

	log.Printf("DataLoader Batch Function finished for keys: %v", symbolNames)
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
