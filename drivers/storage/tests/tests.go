package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	limiter "turbo-dollop"
)

// TestStorageSequentialAccess verify that store works as expected with a sequential access.
func TestStorageSequentialAccess(t *testing.T, store limiter.Storage) {
	is := require.New(t)
	ctx := context.Background()

	limiter := limiter.New(store, limiter.Rate{
		Limit: 3,
		Unit:  time.Minute,
	})

	// Check counter increment.
	{
		for i := 1; i <= 6; i++ {

			if i <= 3 {

				lctx, err := limiter.Peek(ctx, "foo")
				is.NoError(err)
				is.NotZero(lctx)
				is.Equal(int64(3-(i-1)), lctx.Remaining)
				is.False(lctx.Reached)

			}

			lctx, err := limiter.Get(ctx, "foo")
			is.NoError(err)
			is.NotZero(lctx)

			if i <= 3 {

				is.Equal(int64(3), lctx.Limit)
				is.Equal(int64(3-i), lctx.Remaining)
				is.True((lctx.Reset - time.Now().Unix()) <= 60)
				is.False(lctx.Reached)

				lctx, err = limiter.Peek(ctx, "foo")
				is.NoError(err)
				is.Equal(int64(3-i), lctx.Remaining)
				is.False(lctx.Reached)

			} else {

				is.Equal(int64(3), lctx.Limit)
				is.Equal(int64(0), lctx.Remaining)
				is.True((lctx.Reset - time.Now().Unix()) <= 60)
				is.True(lctx.Reached)

			}
		}
	}

	// Check counter reset.
	{
		lctx, err := limiter.Peek(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		is.Equal(int64(3), lctx.Limit)
		is.Equal(int64(0), lctx.Remaining)
		is.True((lctx.Reset - time.Now().Unix()) <= 60)
		is.True(lctx.Reached)

		lctx, err = limiter.Reset(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		is.Equal(int64(3), lctx.Limit)
		is.Equal(int64(3), lctx.Remaining)
		is.True((lctx.Reset - time.Now().Unix()) <= 60)
		is.False(lctx.Reached)

		lctx, err = limiter.Peek(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		is.Equal(int64(3), lctx.Limit)
		is.Equal(int64(3), lctx.Remaining)
		is.True((lctx.Reset - time.Now().Unix()) <= 60)
		is.False(lctx.Reached)

		lctx, err = limiter.Get(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		lctx, err = limiter.Reset(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		is.Equal(int64(3), lctx.Limit)
		is.Equal(int64(3), lctx.Remaining)
		is.True((lctx.Reset - time.Now().Unix()) <= 60)
		is.False(lctx.Reached)

		lctx, err = limiter.Reset(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		is.Equal(int64(3), lctx.Limit)
		is.Equal(int64(3), lctx.Remaining)
		is.True((lctx.Reset - time.Now().Unix()) <= 60)
		is.False(lctx.Reached)
	}
}

// TestStorageConcurrentAccess verify that store works as expected with a concurrent access.
func TestStorageConcurrentAccess(t *testing.T, store limiter.Storage) {
	is := require.New(t)
	ctx := context.Background()

	limiter := limiter.New(store, limiter.Rate{
		Limit: 100000,
		Unit:  10 * time.Second,
	})

	goroutines := 500
	ops := 500

	wg := &sync.WaitGroup{}
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			for j := 0; j < ops; j++ {
				lctx, err := limiter.Get(ctx, "foo")
				is.NoError(err)
				is.NotZero(lctx)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// BenchmarkStorageSequentialAccess executes a benchmark against a store without parallel setting.
func BenchmarkStorageSequentialAccess(b *testing.B, store limiter.Storage) {
	ctx := context.Background()

	instance := limiter.New(store, limiter.Rate{
		Limit: 100000,
		Unit:  10 * time.Second,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = instance.Get(ctx, "foo")
	}
}

// BenchmarkStorageConcurrentAccess executes a benchmark against a store with parallel setting.
func BenchmarkStorageConcurrentAccess(b *testing.B, store limiter.Storage) {
	ctx := context.Background()

	instance := limiter.New(store, limiter.Rate{
		Limit: 100000,
		Unit:  10 * time.Second,
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = instance.Get(ctx, "foo")
		}
	})
}
