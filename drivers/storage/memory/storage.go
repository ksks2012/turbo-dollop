package memory

import (
	"context"
	"strings"
	"time"

	limiter "turbo-dollop"

	"turbo-dollop/drivers/storage/common"
)

// Storage is the in-memory store.
type Storage struct {
	// Prefix used for the key.
	Prefix string
	// cache used to store values in-memory.
	cache *CacheWrapper
}

// NewStorage creates a new instance of memory store with defaults.
func NewStorage() limiter.Storage {
	return NewStorageWithOptions(limiter.StorageOptions{
		Prefix:  limiter.DefaultPrefix,
		CleanUp: limiter.DefaultCleanUpInterval,
	})
}

// NewStorageWithOptions creates a new instance of memory store with options.
func NewStorageWithOptions(options limiter.StorageOptions) limiter.Storage {
	return &Storage{
		Prefix: options.Prefix,
		cache:  NewCache(options.CleanUp),
	}
}

// Get returns the limit for given identifier.
func (store *Storage) Get(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	count, expiration := store.cache.Increment(store.getCacheKey(key), 1, rate.Unit)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, count)
	return lctx, nil
}

// Increment increments the limit by given count & returns the new limit value for given identifier.
func (store *Storage) Increment(ctx context.Context, key string, count int64, rate limiter.Rate) (limiter.Context, error) {
	newCount, expiration := store.cache.Increment(store.getCacheKey(key), count, rate.Unit)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, newCount)
	return lctx, nil
}

// Peek returns the limit for given identifier, without modification on current values.
func (store *Storage) Peek(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	count, expiration := store.cache.Get(store.getCacheKey(key), rate.Unit)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, count)
	return lctx, nil
}

// Reset returns the limit for given identifier.
func (store *Storage) Reset(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	count, expiration := store.cache.Reset(store.getCacheKey(key), rate.Unit)

	lctx := common.GetContextFromState(time.Now(), rate, expiration, count)
	return lctx, nil
}

// getCacheKey returns the full path for an identifier.
func (store *Storage) getCacheKey(key string) string {
	buffer := strings.Builder{}
	buffer.WriteString(store.Prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}
