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

func TestMemoryStorageConcurrentAccess(t *testing.T) {
	tests.TestStorageConcurrentAccess(t, memory.NewStorageWithOptions(limiter.StorageOptions{
		Prefix:  "limiter:memory:concurrent-test",
		CleanUp: 1 * time.Nanosecond,
	}))
}

func BenchmarkMemoryStorageSequentialAccess(b *testing.B) {
	tests.BenchmarkStorageSequentialAccess(b, memory.NewStorageWithOptions(limiter.StorageOptions{
		Prefix:  "limiter:memory:sequential-benchmark",
		CleanUp: 1 * time.Hour,
	}))
}

func BenchmarkMemoryStorageConcurrentAccess(b *testing.B) {
	tests.BenchmarkStorageConcurrentAccess(b, memory.NewStorageWithOptions(limiter.StorageOptions{
		Prefix:  "limiter:memory:concurrent-benchmark",
		CleanUp: 1 * time.Hour,
	}))
}
