// Package groupby provides event grouping by a configurable key function,
// collecting events into named buckets for downstream batch processing.
package groupby

import "github.com/user/portwatch/internal/alert"

// KeyFunc extracts a grouping key from an event.
type KeyFunc func(alert.Event) string

// Group holds a named collection of events.
type Group struct {
	Key    string
	Events []alert.Event
}

// Grouper partitions a slice of events into named buckets.
type Grouper struct {
	keyFn KeyFunc
}

// New returns a Grouper that partitions events using keyFn.
// keyFn must not be nil.
func New(keyFn KeyFunc) (*Grouper, error) {
	if keyFn == nil {
		return nil, ErrNilKeyFunc
	}
	return &Grouper{keyFn: keyFn}, nil
}

// Apply partitions events into groups. Order within each group is preserved;
// group order follows first-seen key insertion order.
func (g *Grouper) Apply(events []alert.Event) []Group {
	if len(events) == 0 {
		return nil
	}

	index := make(map[string]int)
	var groups []Group

	for _, ev := range events {
		k := g.keyFn(ev)
		if i, ok := index[k]; ok {
			groups[i].Events = append(groups[i].Events, ev)
		} else {
			index[k] = len(groups)
			groups = append(groups, Group{Key: k, Events: []alert.Event{ev}})
		}
	}

	return groups
}

// ByAction is a pre-built KeyFunc that groups events by their Action field.
func ByAction(ev alert.Event) string { return ev.Action }

// ByProtocol is a pre-built KeyFunc that groups events by port protocol.
func ByProtocol(ev alert.Event) string { return ev.Port.Protocol }
