package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	defaultTTL      = 5 * time.Minute
	cleanupInterval = 10 * time.Minute
)

// sharedCache holds the single instance of our memory cache.
var sharedCache *gocache.Cache

// init initializes the shared cache when the package is first used.
func init() {
	sharedCache = gocache.New(defaultTTL, cleanupInterval)
}

// Set adds an item to the cache, replacing any existing item.
// It uses the default cache TTL.
func Set(key string, value *time.Time) {
	if value == nil { // Avoid caching nil pointers explicitly, though go-cache might handle it
		return
	}
	sharedCache.Set(key, value, gocache.DefaultExpiration)
}

// Get retrieves an item from the cache.
// It returns the item or nil, and a bool indicating whether the key was found.
func Get(key string) (*time.Time, bool) {
	val, found := sharedCache.Get(key)
	if !found {
		return nil, false
	}

	// Type assertion to ensure we return the correct type
	cachedTime, ok := val.(*time.Time)
	if !ok {
		// Item found but is not the expected type, treat as not found
		return nil, false
	}
	return cachedTime, true
}
