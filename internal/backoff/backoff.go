// Package backoff provides exponential back-off for transient scan errors.
// It is used by the daemon to avoid hammering the OS when the scanner
// returns repeated failures.
package backoff

import (
	"math"
	"sync"
	"time"
)

// Backoff tracks consecutive failures and returns the next wait duration
// using an exponential strategy capped at MaxDelay.
type Backoff struct {
	mu       sync.Mutex
	failures int

	BaseDelay time.Duration
	MaxDelay  time.Duration
	Multiplier float64
}

// New returns a Backoff with sensible defaults.
func New() *Backoff {
	return &Backoff{
		BaseDelay:  250 * time.Millisecond,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
	}
}

// Failure records one failure and returns the duration the caller should
// wait before retrying.
func (b *Backoff) Failure() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures++
	delay := float64(b.BaseDelay) * math.Pow(b.Multiplier, float64(b.failures-1))
	if delay > float64(b.MaxDelay) {
		delay = float64(b.MaxDelay)
	}
	return time.Duration(delay)
}

// Reset clears the failure counter. Call this after a successful operation.
func (b *Backoff) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
}

// Failures returns the current consecutive failure count.
func (b *Backoff) Failures() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures
}
