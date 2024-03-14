package common

import (
	"time"

	turbodollop "turbo-dollop"
)

// GetContextFromState generate a new limiter.Context from given state.
func GetContextFromState(now time.Time, rate turbodollop.Rate, expiration time.Time, count int64) turbodollop.Context {
	limit := rate.Limit
	remaining := int64(0)
	reached := true

	if count <= limit {
		remaining = limit - count
		reached = false
	}

	return turbodollop.Context{
		Limit:     limit,
		Remaining: remaining,
		Reached:   reached,
	}
}
