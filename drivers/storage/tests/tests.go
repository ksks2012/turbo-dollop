package tests

import (
	"context"
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
