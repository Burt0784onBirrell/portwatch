package cooldown

import (
	"fmt"

	"github.com/dkrichards86/portwatch/internal/alert"
)

// KeyForEvent builds a stable string key from an alert event that
// uniquely identifies the port/protocol/action triple.
func KeyForEvent(e alert.Event) string {
	return fmt.Sprintf("%s:%d:%s", e.Port.Protocol, e.Port.Port, e.Action)
}

// FilterEvents returns only those events that are permitted by the
// Tracker. Events whose key is still within the quiet period are
// silently dropped.
func FilterEvents(t *Tracker, events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if t.Allow(KeyForEvent(e)) {
			out = append(out, e)
		}
	}
	return out
}
