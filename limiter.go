package turbodollop

import "context"

type Context struct {
	Limit     int64
	Remaining int64
	Reset     int64
	Reached   bool
}

type Limiter struct {
	Storage Storage
	Rate    Rate
	Options Options
}

// New returns an instance of Limiter.
func New(storage Storage, rate Rate, options ...Option) *Limiter {
	opt := Options{
		IPv4Mask:           DefaultIPv4Mask,
		IPv6Mask:           DefaultIPv6Mask,
		TrustForwardHeader: false,
	}
	for _, o := range options {
		o(&opt)
	}
	return &Limiter{
		Storage: storage,
		Rate:    rate,
		Options: opt,
	}
}

func NewLimiter(storage Storage, rate Rate, options ...Option) *Limiter {
	opt := Options{
		IPv4Mask:           DefaultIPv4Mask,
		IPv6Mask:           DefaultIPv6Mask,
		TrustForwardHeader: false,
	}
	for _, o := range options {
		o(&opt)
	}
	return &Limiter{
		Storage: storage,
		Rate:    rate,
		Options: opt,
	}
}

func (limiter *Limiter) Get(ctx context.Context, key string) (Context, error) {
	return limiter.Storage.Get(ctx, key, limiter.Rate)
}

func (limiter *Limiter) Peek(ctx context.Context, key string) (Context, error) {
	return limiter.Storage.Peek(ctx, key, limiter.Rate)
}

func (limiter *Limiter) Reset(ctx context.Context, key string) (Context, error) {
	return limiter.Storage.Reset(ctx, key, limiter.Rate)
}

func (limiter *Limiter) Increment(ctx context.Context, key string, count int64) (Context, error) {
	return limiter.Storage.Increment(ctx, key, count, limiter.Rate)
}
