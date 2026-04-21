// Package baseline provides a mechanism for capturing and comparing a
// "known-good" port state. Deviations from the baseline are surfaced as
// events so operators can distinguish expected from unexpected changes.
package baseline

import (
	"sync"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/scanner"
)

// Baseline holds a frozen snapshot of ports that are considered normal.
type Baseline struct {
	mu          sync.RWMutex
	ports       scanner.PortSet
	capturedAt  time.Time
}

// New returns an empty Baseline. Call Capture to populate it.
func New() *Baseline {
	return &Baseline{}
}

// Capture replaces the current baseline with the supplied PortSet.
func (b *Baseline) Capture(ports scanner.PortSet) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ports = ports.Clone()
	b.capturedAt = time.Now()
}

// CapturedAt returns the time the baseline was last captured.
// The zero value is returned if no baseline has been captured yet.
func (b *Baseline) CapturedAt() time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.capturedAt
}

// IsSet reports whether a baseline has been captured.
func (b *Baseline) IsSet() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return !b.capturedAt.IsZero()
}

// Deviation compares current against the baseline and returns alert events
// for any ports that have opened or closed relative to it.
// If no baseline has been captured, Deviation returns nil.
func (b *Baseline) Deviation(current scanner.PortSet) []alert.Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.capturedAt.IsZero() {
		return nil
	}
	diff := scanner.Compare(b.ports, current)
	return alert.BuildEvents(diff)
}

// Snapshot returns a copy of the current baseline PortSet.
func (b *Baseline) Snapshot() scanner.PortSet {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ports.Clone()
}
