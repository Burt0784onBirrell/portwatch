// Package limiter provides a token-bucket style rate limiter that caps the
// number of alert events dispatched per unit of time across all notifiers.
package limiter

import (
	"sync"
	"time"
)

// Clock allows the wall-clock to be replaced in tests.
type Clock func() time.Time

// Limiter tracks how many events have been emitted within a rolling window and
// drops events once the configured burst limit is reached.
type Limiter struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	clock   Clock
	buckets []time.Time
}

// New creates a Limiter that allows at most max events per window.
func New(max int, window time.Duration) *Limiter {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, clock Clock) *Limiter {
	if max < 1 {
		max = 1
	}
	return &Limiter{
		max:    max,
		window: window,
		clock:  clock,
	}
}

// Allow returns true when the event should be forwarded and false when the
// burst limit for the current window has been exceeded.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	cutoff := now.Add(-l.window)

	// Evict timestamps outside the window.
	valid := l.buckets[:0]
	for _, t := range l.buckets {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	l.buckets = valid

	if len(l.buckets) >= l.max {
		return false
	}
	l.buckets = append(l.buckets, now)
	return true
}

// Remaining returns the number of events that can still be emitted in the
// current window without being dropped.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets {
		if t.After(cutoff) {
			count++
		}
	}
	r := l.max - count
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded timestamps, immediately restoring the full burst
// budget.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = l.buckets[:0]
}
