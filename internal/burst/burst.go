// Package burst detects sudden spikes in port-change events over a
// rolling window and suppresses or flags them as burst activity.
package burst

import (
	"sync"
	"time"
)

// Detector tracks event counts in a sliding window and reports whether
// the current rate constitutes a burst.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	buckets   []entry
	now       func() time.Time
}

type entry struct {
	at    time.Time
	count int
}

// New returns a Detector that considers activity a burst when more than
// threshold events are recorded within the given window.
func New(window time.Duration, threshold int) *Detector {
	return newWithClock(window, threshold, time.Now)
}

func newWithClock(window time.Duration, threshold int, now func() time.Time) *Detector {
	if threshold < 1 {
		threshold = 1
	}
	return &Detector{
		window:    window,
		threshold: threshold,
		now:       now,
	}
}

// Record adds n events at the current instant.
func (d *Detector) Record(n int) {
	if n <= 0 {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	d.buckets = append(d.buckets, entry{at: d.now(), count: n})
}

// IsBurst reports whether the total count within the window exceeds the
// configured threshold.
func (d *Detector) IsBurst() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	return d.total() > d.threshold
}

// Total returns the event count currently within the window.
func (d *Detector) Total() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	return d.total()
}

// Reset clears all recorded entries.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buckets = d.buckets[:0]
}

func (d *Detector) evict() {
	cutoff := d.now().Add(-d.window)
	i := 0
	for i < len(d.buckets) && d.buckets[i].at.Before(cutoff) {
		i++
	}
	d.buckets = d.buckets[i:]
}

func (d *Detector) total() int {
	sum := 0
	for _, e := range d.buckets {
		sum += e.count
	}
	return sum
}
