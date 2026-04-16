// Package truncate provides a middleware that limits the number of events
// forwarded to downstream notifiers in a single dispatch cycle.
//
// When the number of events exceeds the configured cap the slice is trimmed to
// the cap and a synthetic summary event is appended so operators know that
// events were dropped.
package truncate

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// Truncator caps the number of alert events per dispatch cycle.
type Truncator struct {
	max int
}

// New returns a Truncator that allows at most max events through.
// If max is less than 1 it is clamped to 1.
func New(max int) *Truncator {
	if max < 1 {
		max = 1
	}
	return &Truncator{max: max}
}

// Apply returns a (possibly truncated) copy of events.
// When truncation occurs the last element of the returned slice is replaced
// with a synthetic "truncated" notice so the total length never exceeds max.
func (t *Truncator) Apply(events []alert.Event) []alert.Event {
	if len(events) <= t.max {
		out := make([]alert.Event, len(events))
		copy(out, events)
		return out
	}

	// Keep max-1 real events and add one summary notice.
	keep := t.max - 1
	if keep < 0 {
		keep = 0
	}

	out := make([]alert.Event, keep, t.max)
	copy(out, events[:keep])

	dropped := len(events) - keep
	notice := events[0] // copy header fields from first event
	notice.Process = fmt.Sprintf("[truncated: %d events dropped]", dropped)
	out = append(out, notice)
	return out
}
