package trap

import "context"

// Trap defines the interface that all honeypot services must implement.
type Trap interface {
	// Start starts the trap service. It should block until the service stops or fails.
	// The context can be used to signal cancellation.
	Start(ctx context.Context) error

	// Shutdown gracefully shuts down the trap service.
	Shutdown(ctx context.Context) error
}
