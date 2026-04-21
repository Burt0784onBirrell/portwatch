// Package trend tracks the rate of port change events over time,
// providing a rolling view of how frequently ports are opening or closing.
package trend

import (
	"sync"
	"time"
)

// Direction represents whether the trend is rising, falling, or stable.
type Direction string

const (
	Rising  Direction = "rising"
	Falling Direction = "falling"
	Stable  Direction = "stable"
)

// Point is a single observation recorded at a point in time.
type Point struct {
	At    time.Time
	Count int
}

// Tracker accumulates event counts within a sliding window and exposes
// the current trend direction by comparing recent halves of the window.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	points  []Point
	clock   func() time.Time
}

// New returns a Tracker with the given sliding-window duration.
func New(window time.Duration) *Tracker {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Tracker {
	return &Tracker{window: window, clock: clock}
}

// Record adds n events to the current observation bucket.
func (t *Tracker) Record(n int) {
	if n <= 0 {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	t.evict(now)
	t.points = append(t.points, Point{At: now, Count: n})
}

// Direction returns Rising, Falling, or Stable based on whether the
// second half of the window has more, fewer, or equal events than the first.
func (t *Tracker) Direction() Direction {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	t.evict(now)
	if len(t.points) == 0 {
		return Stable
	}
	mid := now.Add(-t.window / 2)
	var first, second int
	for _, p := range t.points {
		if p.At.Before(mid) {
			first += p.Count
		} else {
			second += p.Count
		}
	}
	switch {
	case second > first:
		return Rising
	case second < first:
		return Falling
	default:
		return Stable
	}
}

// Total returns the sum of all event counts within the current window.
func (t *Tracker) Total() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict(t.clock())
	var sum int
	for _, p := range t.points {
		sum += p.Count
	}
	return sum
}

// evict removes observations that have fallen outside the window.
// Caller must hold t.mu.
func (t *Tracker) evict(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.points) && t.points[i].At.Before(cutoff) {
		i++
	}
	t.points = t.points[i:]
}
