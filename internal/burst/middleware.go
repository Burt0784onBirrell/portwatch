package burst

import "github.com/joshbeard/portwatch/internal/alert"

// FilterWhenBursting returns events unchanged when activity is below the
// burst threshold. When a burst is detected it returns a nil slice so that
// downstream notifiers are not overwhelmed.
//
// The detector is updated with the number of events in every call.
func FilterWhenBursting(d *Detector, events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	d.Record(len(events))
	if d.IsBurst() {
		return nil
	}
	return events
}
