// Package pressure tracks the rate of port-change events over a sliding
// window and exposes a simple Load value in the range [0.0, 1.0].  When Load
// exceeds a configurable threshold the detector signals high pressure, which
// callers can use to shed work (e.g. skip expensive enrichment steps).
package pressure

import (
	"sync"
	"time"
)

// Detector measures event throughput and reports whether the system is under
// high pressure.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	capacity  int
	threshold float64
	buckets   []time.Time
	now       func() time.Time
}

// New returns a Detector that considers capacity events within window to be
// full load.  threshold is the fraction of capacity [0,1] above which
// IsHigh returns true.
func New(window time.Duration, capacity int, threshold float64) *Detector {
	if threshold < 0 {
		threshold = 0
	}
	if threshold > 1 {
		threshold = 1
	}
	if capacity < 1 {
		capacity = 1
	}
	return newWithClock(window, capacity, threshold, time.Now)
}

func newWithClock(window time.Duration, capacity int, threshold float64, now func() time.Time) *Detector {
	return &Detector{
		window:    window,
		capacity:  capacity,
		threshold: threshold,
		now:       now,
	}
}

// Record registers n events occurring now.
func (d *Detector) Record(n int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	d.evict(now)
	for i := 0; i < n; i++ {
		d.buckets = append(d.buckets, now)
	}
}

// Load returns the fraction of capacity consumed in the current window.
func (d *Detector) Load() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict(d.now())
	return float64(len(d.buckets)) / float64(d.capacity)
}

// IsHigh returns true when Load exceeds the configured threshold.
func (d *Detector) IsHigh() bool {
	return d.Load() >= d.threshold
}

// evict removes timestamps that have fallen outside the window. Must be called
// with d.mu held.
func (d *Detector) evict(now time.Time) {
	cutoff := now.Add(-d.window)
	i := 0
	for i < len(d.buckets) && d.buckets[i].Before(cutoff) {
		i++
	}
	d.buckets = d.buckets[i:]
}
