// Package label provides a simple key/value tagging layer that attaches
// static metadata labels to alert events before dispatch.
package label

import (
	"errors"

	"github.com/user/portwatch/internal/alert"
)

// Labeler attaches a fixed set of key/value labels to every event it processes.
type Labeler struct {
	labels map[string]string
}

// New returns a Labeler that will stamp each event with the provided labels.
// Returns an error if labels is nil or any key is empty.
func New(labels map[string]string) (*Labeler, error) {
	if len(labels) == 0 {
		return nil, errors.New("label: at least one label is required")
	}
	for k := range labels {
		if k == "" {
			return nil, errors.New("label: empty key is not allowed")
		}
	}
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return &Labeler{labels: copy}, nil
}

// Apply returns a new slice of events where each event's Meta map is augmented
// with the Labeler's labels. Existing keys on the event are not overwritten.
func (l *Labeler) Apply(events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, len(events))
	for i, ev := range events {
		merged := make(map[string]string, len(l.labels))
		for k, v := range l.labels {
			merged[k] = v
		}
		for k, v := range ev.Meta {
			merged[k] = v // event labels win
		}
		ev.Meta = merged
		out[i] = ev
	}
	return out
}
