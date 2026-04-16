// Package window provides a sliding-window event counter that tracks how
// many events occurred within a rolling time duration.
package window

import (
	"sync"
	"time"
)

// Clock allows tests to inject a fake time source.
type Clock func() time.Time

// Counter is a thread-safe sliding-window counter.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	clock    Clock
	buckets  []time.Time
}

// New returns a Counter that tracks events within the given window duration.
func New(window time.Duration) *Counter {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock Clock) *Counter {
	return &Counter{window: window, clock: clock}
}

// Record registers one event at the current time.
func (c *Counter) Record() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buckets = append(c.buckets, c.clock())
	c.evict()
}

// Count returns the number of events recorded within the current window.
func (c *Counter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict()
	return len(c.buckets)
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buckets = nil
}

// evict removes entries that have fallen outside the window. Must be called
// with c.mu held.
func (c *Counter) evict() {
	cutoff := c.clock().Add(-c.window)
	i := 0
	for i < len(c.buckets) && c.buckets[i].Before(cutoff) {
		i++
	}
	c.buckets = c.buckets[i:]
}
