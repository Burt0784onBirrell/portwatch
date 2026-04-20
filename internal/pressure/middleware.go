package pressure

import "github.com/danvolchek/portwatch/internal/alert"

// FilterWhenHigh returns events unchanged when the detector reports normal
// load.  When pressure is high it returns an empty slice, effectively shedding
// the batch so that downstream notifiers are not overwhelmed.
//
// Record is called with the batch size before the load check so that the
// detector always has an accurate view of incoming throughput.
func FilterWhenHigh(d *Detector, events []alert.Event) []alert.Event {
	d.Record(len(events))
	if d.IsHigh() {
		return []alert.Event{}
	}
	return events
}
