package census

type Config struct {
	// The smallest desired cluster size (in consumers)
	MinimumClusterSize int
}

type Client interface {
	// Sync updates the backend with information on our node's population.
	Sync()

	// AutoSync automatically calls Sync() periodically. Frequency should be set
	// automatically via a bounded backoff algorithm.
	AutoSync()

	// GetNeighbors gets all cluster neighbors, excluding itself in the output.
	GetNeighbors()
}
