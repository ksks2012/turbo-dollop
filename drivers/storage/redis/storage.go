package redis

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"turbo-dollop/drivers/storage/common"

	turbodollop "turbo-dollop"

	libredis "github.com/redis/go-redis/v9"
)

// Redis client interface
type Client interface {
	Get(ctx context.Context, key string) *libredis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *libredis.StatusCmd
	Watch(ctx context.Context, handler func(*libredis.Tx) error, keys ...string) error
	Del(ctx context.Context, keys ...string) *libredis.IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *libredis.BoolCmd
	EvalSha(ctx context.Context, sha string, keys []string, args ...interface{}) *libredis.Cmd
	ScriptLoad(ctx context.Context, script string) *libredis.StringCmd
	Incr(ctx context.Context, key string) *libredis.IntCmd
	IncrBy(ctx context.Context, key string, value int64) *libredis.IntCmd
}

func NewStorage(client Client) (turbodollop.Storage, error) {
	return NewStorageWithOptions(client, turbodollop.StorageOptions{
		Prefix:   turbodollop.DefaultPrefix,
		CleanUp:  turbodollop.DefaultCleanUpInterval,
		MaxRetry: turbodollop.DefaultMaxRetry,
	})
}

type Storage struct {
	client   Client
	Prefix   string
	MaxRetry int
	rwMutex  sync.RWMutex
}

func NewStorageWithOptions(client Client, options turbodollop.StorageOptions) (turbodollop.Storage, error) {
	storage := &Storage{
		client:   client,
		Prefix:   options.Prefix,
		MaxRetry: options.MaxRetry,
	}

	return storage, nil
}

func (storage *Storage) Get(ctx context.Context, key string, rate turbodollop.Rate) (turbodollop.Context, error) {
	storage.rwMutex.RLock()
	defer storage.rwMutex.RUnlock()
	val, err := storage.client.Incr(ctx, storage.getCacheKey(key)).Result()
	if err != nil {
		return currentContext(0, rate)
	}
	return currentContext(int64(val), rate)
}

func (storage *Storage) Peek(ctx context.Context, key string, rate turbodollop.Rate) (turbodollop.Context, error) {
	storage.rwMutex.RLock()
	defer storage.rwMutex.RUnlock()
	val, err := storage.client.Get(ctx, storage.getCacheKey(key)).Result()
	if err == libredis.Nil {
		// Create the key if it does not exist
		_, err := storage.client.SetNX(ctx, storage.getCacheKey(key), 0, 0).Result()
		if err != nil {
			return turbodollop.Context{}, err
		}
		return currentContext(0, rate)
	} else if err != nil {
		return turbodollop.Context{}, err
	}
	numVal, numErr := strconv.ParseInt(val, 10, 64)
	if numErr != nil {
		panic(err)
	}
	return currentContext(int64(numVal), rate)
}

func (storage *Storage) Increment(ctx context.Context, key string, count int64, rate turbodollop.Rate) (turbodollop.Context, error) {
	storage.rwMutex.RLock()
	defer storage.rwMutex.RUnlock()
	// TODO:
	val, err := storage.client.IncrBy(ctx, storage.getCacheKey(key), count).Result()
	if err != nil {
		return turbodollop.Context{}, err
	}
	return currentContext(val, rate)
}

// Reset returns the limit for given identifier which is set to zero.
func (storage *Storage) Reset(ctx context.Context, key string, rate turbodollop.Rate) (turbodollop.Context, error) {
	storage.rwMutex.RLock()
	defer storage.rwMutex.RUnlock()
	_, err := storage.client.Set(ctx, storage.getCacheKey(key), 0, 0).Result()
	if err != nil {
		return turbodollop.Context{}, err
	}
	return currentContext(0, rate)
}

// getCacheKey returns the full path for an identifier.
func (storage *Storage) getCacheKey(key string) string {
	buffer := strings.Builder{}
	buffer.WriteString(storage.Prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}

func currentContext(val int64, rate turbodollop.Rate) (turbodollop.Context, error) {
	now := time.Now()
	expiration := now.Add(rate.Unit)

	return common.GetContextFromState(now, rate, expiration, val), nil
}
