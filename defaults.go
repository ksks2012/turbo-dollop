package turbodollop

import "time"

const (
	// DefaultPrefix is the default prefix to use for the key in the store.
	DefaultPrefix = "turbodollop"

	// DefaultMaxRetry is the default maximum number of key retries under
	// race condition (mainly used with database-based stores).
	DefaultMaxRetry = 3

	// DefaultCleanUpInterval is the default time duration for cleanup.
	DefaultCleanUpInterval = 30 * time.Second

	// Default Setting for Token Bucket
	DefaultRefillRate  = 1 * time.Second
	DefaultTokenBucket = 5

	// Default Setting for Leaky Bucket
	DefaultOutflowRate = 1 * time.Second
	DefaultQueueSize   = 5

	// Default Setting for Sliding Window
	DefaultOutdateing = 1 * time.Second
	DefaultMaxLog     = 5
)
