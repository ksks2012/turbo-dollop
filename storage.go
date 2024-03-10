package turbodollop

import (
	"context"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string, rate Rate) (Context, error)
	Reset(ctx context.Context, key string, rate Rate) (Context, error)
	Increment(ctx context.Context, key string, count int64, rate Rate) (Context, error)
}

type StorageOptions struct {
	Prefix   string
	MaxRetry int
	CleanUp  time.Duration
}
