// Package quota enforces per-key event quotas over a rolling time window.
// When a key exceeds its configured maximum, further events are dropped until
// the window resets.
package quota

import (
	"sync"
	"time"
)

// clock allows tests to inject a fake time source.
type clock func() time.Time

// entry tracks usage for a single key within the current window.
type entry struct {
	count     int
	windowEnd time.Time
}

// Quota enforces a maximum number of events per key per window duration.
type Quota struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	entries map[string]*entry
	now     clock
}

// New creates a Quota that allows at most max events per key per window.
// max is clamped to a minimum of 1.
func New(max int, window time.Duration) *Quota {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, now clock) *Quota {
	if max < 1 {
		max = 1
	}
	return &Quota{
		max:     max,
		window:  window,
		entries: make(map[string]*entry),
		now:     now,
	}
}

// Allow returns true if the key has not yet reached its quota for the current
// window. Each call that returns true consumes one unit of quota.
func (q *Quota) Allow(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	e, ok := q.entries[key]
	if !ok || now.After(e.windowEnd) {
		q.entries[key] = &entry{count: 1, windowEnd: now.Add(q.window)}
		return true
	}
	if e.count >= q.max {
		return false
	}
	e.count++
	return true
}

// Reset clears the quota state for key, allowing it to pass immediately.
func (q *Quota) Reset(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.entries, key)
}

// Remaining returns how many more events key may emit in the current window.
// Returns max when no usage has been recorded.
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	e, ok := q.entries[key]
	if !ok || now.After(e.windowEnd) {
		return q.max
	}
	remaining := q.max - e.count
	if remaining < 0 {
		return 0
	}
	return remaining
}
