// Package sequence assigns monotonically increasing sequence numbers to
// alert events, making it easy to detect dropped or out-of-order deliveries.
package sequence

import (
	"fmt"
	"sync/atomic"

	"portwatch/internal/alert"
)

// Sequencer stamps each event with a sequence number via its Tags map.
type Sequencer struct {
	counter uint64
	field   string
}

// New returns a Sequencer that writes the sequence number into the given tag
// field name (e.g. "seq").
func New(field string) (*Sequencer, error) {
	if field == "" {
		return nil, ErrEmptyField
	}
	return &Sequencer{field: field}, nil
}

// Annotate stamps each event with the next sequence number and returns the
// updated slice. The original events are not modified.
func (s *Sequencer) Annotate(events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, len(events))
	for i, e := range events {
		n := atomic.AddUint64(&s.counter, 1)
		if e.Tags == nil {
			e.Tags = make(map[string]string)
		} else {
			copy := make(map[string]string, len(e.Tags))
			for k, v := range e.Tags {
				copy[k] = v
			}
			e.Tags = copy
		}
		e.Tags[s.field] = fmt.Sprintf("%d", n)
		out[i] = e
	}
	return out
}

// Reset sets the internal counter back to zero. Useful in tests.
func (s *Sequencer) Reset() {
	atomic.StoreUint64(&s.counter, 0)
}
