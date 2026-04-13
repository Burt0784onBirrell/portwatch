// Package circuitbreaker implements a simple circuit breaker pattern
// to prevent cascading failures when notifiers or external services are
// repeatedly failing. Once the failure threshold is exceeded the circuit
// opens and calls are rejected until the reset timeout elapses.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and the call is rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	failures     int
	threshold    int
	resetTimeout time.Duration
	openedAt     time.Time
	state        State
	now          func() time.Time
}

// New returns a Breaker that opens after threshold consecutive failures
// and resets after resetTimeout.
func New(threshold int, resetTimeout time.Duration) *Breaker {
	return newWithClock(threshold, resetTimeout, time.Now)
}

func newWithClock(threshold int, resetTimeout time.Duration, now func() time.Time) *Breaker {
	return &Breaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
		now:          now,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit
// is open. Callers must follow up with Record to report the outcome.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateOpen {
		if b.now().Sub(b.openedAt) >= b.resetTimeout {
			// Half-open: allow a single probe through.
			b.state = StateClosed
			b.failures = 0
		} else {
			return ErrOpen
		}
	}
	return nil
}

// Record registers the outcome of a call. A non-nil err counts as a
// failure; a nil err resets the failure counter.
func (b *Breaker) Record(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err == nil {
		b.failures = 0
		return
	}
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current circuit state.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
