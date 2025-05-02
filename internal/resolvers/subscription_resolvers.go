package resolvers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mxcoppell/graphql-resolver-batch-cache/internal/gen/graph/model"
)

// SymbolUpdatesImpl provides the implementation logic for the Subscription.symbolUpdates resolver.
func SymbolUpdatesImpl(ctx context.Context, names []string) (<-chan *model.SymbolDefinition, error) {
	log.Printf("Subscription.symbolUpdates called with %d symbols", len(names))

	// Create a channel to send updates
	ch := make(chan *model.SymbolDefinition, 1)

	// Start a goroutine to send periodic updates
	go func() {
		defer close(ch)

		// Keep track of the current index in the names slice
		index := 0

		// Send an update every 2 seconds until context is cancelled
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Subscription context done, stopping updates.")
				return
			case <-ticker.C:
				// Get the current symbol name
				name := names[index]

				// Create a symbol definition (NextExDividendDate will be resolved downstream)
				desc := fmt.Sprintf("Updated Description for Symbol %s at %s", name, time.Now().Format(time.Kitchen))
				symbol := &model.SymbolDefinition{
					Name:        name,
					Description: &desc,
				}

				// Send it to the channel
				log.Printf("Sending update for symbol %s", name)
				ch <- symbol

				// Move to the next name (round-robin)
				index = (index + 1) % len(names)
			}
		}
	}()

	return ch, nil
}
