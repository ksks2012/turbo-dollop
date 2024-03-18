package redis_test

import (
	"os"
	"testing"

	limiter "turbo-dollop"

	"turbo-dollop/drivers/storage/redis"

	"turbo-dollop/drivers/storage/tests"

	libredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisStorageSequentialAccess(t *testing.T) {
	is := require.New(t)

	client, err := newRedisClient()
	is.NoError(err)
	is.NotNil(client)

	storage, err := redis.NewStorageWithOptions(client, limiter.StorageOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(storage)

	tests.TestStorageSequentialAccess(t, storage)
}

func TestRedisStoreConcurrentAccess(t *testing.T) {
	is := require.New(t)

	client, err := newRedisClient()
	is.NoError(err)
	is.NotNil(client)

	store, err := redis.NewStorageWithOptions(client, limiter.StorageOptions{
		Prefix: "limiter:redis:concurrent-test",
	})
	is.NoError(err)
	is.NotNil(store)

	tests.TestStorageConcurrentAccess(t, store)
}

func newRedisClient() (*libredis.Client, error) {
	uri := "redis://localhost:6379/0"
	if os.Getenv("REDIS_URI") != "" {
		uri = os.Getenv("REDIS_URI")
	}

	opt, err := libredis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	client := libredis.NewClient(opt)
	return client, nil
}
