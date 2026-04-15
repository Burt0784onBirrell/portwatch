package limiter

import "github.com/rnemeth90/portwatch/internal/alert"

// FilterEvents forwards only those events that the Limiter permits, dropping
// the rest. The slice returned is always non-nil.
func FilterEvents(l *Limiter, events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return []alert.Event{}
	}

	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if l.Allow() {
			out = append(out, e)
		}
	}
	return out
}
