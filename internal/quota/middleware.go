package quota

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// KeyForEvent returns a quota key derived from the event's port and protocol.
func KeyForEvent(e alert.Event) string {
	return fmt.Sprintf("%s:%d", e.Port.Protocol, e.Port.Number)
}

// FilterEvents drops any events whose per-key quota has been exhausted and
// returns the remaining allowed events. Order is preserved.
func FilterEvents(q *Quota, events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if q.Allow(KeyForEvent(e)) {
			out = append(out, e)
		}
	}
	return out
}
