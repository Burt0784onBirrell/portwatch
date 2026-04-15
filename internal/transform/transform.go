// Package transform provides a composable event transformation pipeline
// that applies a sequence of mapping functions to alert events.
package transform

import "github.com/user/portwatch/internal/alert"

// MapFunc is a function that transforms a single alert event.
type MapFunc func(alert.Event) alert.Event

// Transformer applies an ordered list of MapFuncs to each event in a slice.
type Transformer struct {
	fns []MapFunc
}

// New creates a Transformer with the provided mapping functions.
// Functions are applied in the order they are supplied.
func New(fns ...MapFunc) *Transformer {
	copy := make([]MapFunc, len(fns))
	for i, f := range fns {
		copy[i] = f
	}
	return &Transformer{fns: copy}
}

// Add appends a new MapFunc to the end of the transformation chain.
func (t *Transformer) Add(fn MapFunc) {
	t.fns = append(t.fns, fn)
}

// Apply runs all registered MapFuncs over each event and returns the
// resulting slice. The original slice is never modified.
func (t *Transformer) Apply(events []alert.Event) []alert.Event {
	if len(events) == 0 || len(t.fns) == 0 {
		return events
	}

	out := make([]alert.Event, len(events))
	for i, e := range events {
		for _, fn := range t.fns {
			e = fn(e)
		}
		out[i] = e
	}
	return out
}

// Compose returns a new Transformer that chains the stages of t followed
// by the stages of other, without mutating either.
func (t *Transformer) Compose(other *Transformer) *Transformer {
	combined := make([]MapFunc, 0, len(t.fns)+len(other.fns))
	combined = append(combined, t.fns...)
	combined = append(combined, other.fns...)
	return &Transformer{fns: combined}
}
