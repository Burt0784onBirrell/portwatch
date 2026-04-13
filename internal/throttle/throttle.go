// Package throttle provides a scan-level throttle that prevents the daemon
// from hammering the OS with port scans when the system is under load.
package throttle

import (
	"context"
	"time"
)

// Scanner is the interface satisfied by any port scanner the throttle wraps.
type Scanner interface {
	Scan(ctx context.Context) (interface{}, error)
}

// ScanFunc is a convenience adapter so ordinary functions can be used as
// Scanner implementations.
type ScanFunc func(ctx context.Context) (interface{}, error)

func (f ScanFunc) Scan(ctx context.Context) (interface{}, error) { return f(ctx) }

// Throttle enforces a minimum interval between successive scans.
type Throttle struct {
	minInterval time.Duration
	last        time.Time
	clock       func() time.Time
}

// New returns a Throttle that will suppress scans that arrive sooner than
// minInterval after the previous one.
func New(minInterval time.Duration) *Throttle {
	return &Throttle{
		minInterval: minInterval,
		clock:       time.Now,
	}
}

// newWithClock is used in tests to inject a fake clock.
func newWithClock(minInterval time.Duration, clock func() time.Time) *Throttle {
	return &Throttle{minInterval: minInterval, clock: clock}
}

// Allow reports whether enough time has elapsed since the last scan.
// When it returns true it also records the current time as the new baseline.
func (t *Throttle) Allow() bool {
	now := t.clock()
	if t.last.IsZero() || now.Sub(t.last) >= t.minInterval {
		t.last = now
		return true
	}
	return false
}

// Reset clears the last-scan timestamp so the next call to Allow always
// succeeds regardless of the configured interval.
func (t *Throttle) Reset() {
	t.last = time.Time{}
}

// Remaining returns how long the caller must wait before Allow will return
// true again. A zero or negative duration means the throttle is ready.
func (t *Throttle) Remaining() time.Duration {
	if t.last.IsZero() {
		return 0
	}
	return t.minInterval - t.clock().Sub(t.last)
}
