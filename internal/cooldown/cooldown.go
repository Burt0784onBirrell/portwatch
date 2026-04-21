// Package cooldown provides a per-key cooldown tracker that prevents
// the same event key from being acted upon more than once within a
// configurable quiet period.
package cooldown

import (
	"sync"
	"time"
)

// clock abstracts time so tests can control it.
type clock func() time.Time

// Tracker records the last activation time for each key and reports
// whether a new activation should be allowed.
type Tracker struct {
	mu       sync.Mutex
	entries  map[string]time.Time
	period   time.Duration
	now      clock
}

// New creates a Tracker with the given quiet period.
func New(period time.Duration) *Tracker {
	return newWithClock(period, time.Now)
}

func newWithClock(period time.Duration, c clock) *Tracker {
	return &Tracker{
		entries: make(map[string]time.Time),
		period:  period,
		now:     c,
	}
}

// Allow returns true if key has not been activated within the quiet
// period, and records the activation time. Returns false otherwise.
func (t *Tracker) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.entries[key]; ok && now.Sub(last) < t.period {
		return false
	}
	t.entries[key] = now
	return true
}

// Reset removes the cooldown entry for key, allowing it to fire
// immediately on the next call to Allow.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Len returns the number of tracked keys.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
