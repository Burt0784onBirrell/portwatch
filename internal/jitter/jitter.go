// Package jitter adds randomised delay to scan intervals to avoid
// thundering-herd problems when multiple portwatch instances run in parallel.
package jitter

import (
	"math/rand"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0, 1).
type Source func() float64

// Jitter applies a random offset to a base duration.
type Jitter struct {
	factor float64
	source Source
}

// New creates a Jitter that spreads delays by up to factor * base.
// factor is clamped to [0, 1].
func New(factor float64) *Jitter {
	return newWithSource(factor, rand.Float64)
}

func newWithSource(factor float64, src Source) *Jitter {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return &Jitter{factor: factor, source: src}
}

// Apply returns base + a random offset in [0, factor*base).
// When factor is 0 the original duration is returned unchanged.
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if j.factor == 0 || base <= 0 {
		return base
	}
	max := float64(base) * j.factor
	offset := time.Duration(j.source() * max)
	return base + offset
}

// Reset returns a new ticker channel that fires after Apply(interval).
// Callers should stop the returned ticker when done.
func (j *Jitter) Reset(interval time.Duration) *time.Timer {
	return time.NewTimer(j.Apply(interval))
}
