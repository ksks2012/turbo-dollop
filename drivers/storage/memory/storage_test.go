package memory_test

import (
	"testing"
	"time"

	limiter "turbo-dollop"
	"turbo-dollop/drivers/storage/memory"
	"turbo-dollop/drivers/storage/tests"
)

func TestMemoryStorageSequentialAccess(t *testing.T) {
	tests.TestStorageSequentialAccess(t, memory.NewStorageWithOptions(limiter.StorageOptions{
		Prefix:  "limiter:memory:sequential-test",
		CleanUp: 30 * time.Second,
	}))
}