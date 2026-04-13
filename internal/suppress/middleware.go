package suppress

import (
	"fmt"

	"github.com/joshbeard/portwatch/internal/alert"
)

// KeyForEvent returns a canonical suppression key for an alert event.
func KeyForEvent(e alert.Event) string {
	return fmt.Sprintf("%s:%d:%s", e.Port.Protocol, e.Port.Number, e.Action)
}

// FilterEvents returns only those events that are not currently suppressed.
// Events that pass through are NOT automatically suppressed; callers must
// call Store.Suppress explicitly if desired.
func FilterEvents(store *Store, events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if !store.IsSuppressed(KeyForEvent(e)) {
			out = append(out, e)
		}
	}
	return out
}
