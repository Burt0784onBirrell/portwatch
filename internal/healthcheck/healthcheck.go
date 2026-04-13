package healthcheck

import (
	"sync"
	"time"
)

// Status represents the current health of the daemon.
type Status struct {
	Healthy     bool      `json:"healthy"`
	LastScan    time.Time `json:"last_scan"`
	ScanCount   int64     `json:"scan_count"`
	ErrorCount  int64     `json:"error_count"`
	UptimeSince time.Time `json:"uptime_since"`
}

// Checker tracks daemon liveness and exposes a health status.
type Checker struct {
	mu          sync.RWMutex
	lastScan    time.Time
	scanCount   int64
	errorCount  int64
	uptimeSince time.Time
	maxStaleness time.Duration
}

// New creates a Checker with the given staleness threshold.
// If the last scan is older than maxStaleness the checker reports unhealthy.
func New(maxStaleness time.Duration) *Checker {
	return &Checker{
		uptimeSince:  time.Now(),
		maxStaleness: maxStaleness,
	}
}

// RecordScan marks a successful scan tick.
func (c *Checker) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastScan = time.Now()
	c.scanCount++
}

// RecordError increments the error counter without updating the scan time.
func (c *Checker) RecordError() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorCount++
}

// Status returns a point-in-time snapshot of health.
func (c *Checker) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	healthy := !c.lastScan.IsZero() &&
		time.Since(c.lastScan) <= c.maxStaleness

	return Status{
		Healthy:     healthy,
		LastScan:    c.lastScan,
		ScanCount:   c.scanCount,
		ErrorCount:  c.errorCount,
		UptimeSince: c.uptimeSince,
	}
}
