package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of daemon metrics.
type Snapshot struct {
	ScansTotal    int64
	AlertsTotal   int64
	LastScanAt    time.Time
	LastChangeAt  time.Time
	UptimeSeconds float64
	startedAt     time.Time
}

// Collector accumulates runtime metrics for the daemon.
type Collector struct {
	mu           sync.RWMutex
	scansTotal   int64
	alertsTotal  int64
	lastScanAt   time.Time
	lastChangeAt time.Time
	startedAt    time.Time
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{startedAt: time.Now()}
}

// RecordScan increments the scan counter and updates the last-scan timestamp.
func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
}

// RecordChange increments the alert counter and updates the last-change timestamp.
func (c *Collector) RecordChange(count int) {
	if count == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal += int64(count)
	c.lastChangeAt = time.Now()
}

// Snapshot returns a consistent copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:    c.scansTotal,
		AlertsTotal:   c.alertsTotal,
		LastScanAt:    c.lastScanAt,
		LastChangeAt:  c.lastChangeAt,
		UptimeSeconds: time.Since(c.startedAt).Seconds(),
		startedAt:     c.startedAt,
	}
}
