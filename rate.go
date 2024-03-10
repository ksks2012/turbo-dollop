package turbodollop

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Rate struct {
	Options string
	Unit    time.Duration
	Limit   int64
}

func NewRateCreator(options string) (Rate, error) {
	rate := Rate{}

	// option include limit and unit like 10|H
	option := strings.Split(options, "|")
	if len(option) != 2 {
		return rate, errors.Errorf("incorrect format '%s'", options)
	}

	units := map[string]time.Duration{
		"MS": time.Millisecond,
		"S":  time.Second,
		"M":  time.Minute,
		"H":  time.Hour,
		"D":  time.Hour * 24,
		"W":  time.Hour * 24 * 7,
	}

	limit, unit := option[0], strings.ToUpper(option[1])

	u, unit_err := units[unit]
	if !unit_err {
		return rate, errors.Errorf("incorrect unit '%s'", unit)
	}

	l, limit_err := strconv.ParseInt(limit, 10, 64)
	if limit_err != nil {
		return rate, errors.Errorf("incorrect limit '%s'", limit)
	}

	rate = Rate{
		Options: options,
		Unit:    u,
		Limit:   l,
	}

	return rate, nil
}
